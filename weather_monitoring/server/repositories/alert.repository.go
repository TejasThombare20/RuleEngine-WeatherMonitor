package repositories

import (
	"database/sql"
	"log"
)

type AlertRepository struct {
	db *sql.DB
}

func NewAlertRepository(db *sql.DB) (*AlertRepository, error) {

	repo := &AlertRepository{db: db}

	if err := repo.initialize(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *AlertRepository) initialize() error {

	_, err := r.db.Exec(`CREATE TABLE IF NOT EXISTS users (
    		id SERIAL PRIMARY KEY,
    		email TEXT NOT NULL UNIQUE,
    		temperature_unit TEXT NOT NULL DEFAULT 'celsius' 
		);`)

	if err != nil {
		log.Println("error while creating a users table")
		return err
	}

	_, err = r.db.Exec(`CREATE TABLE IF NOT EXISTS user_city_thresholds (
    			id SERIAL PRIMARY KEY,
    			user_id INTEGER REFERENCES users(id),
    			city_name TEXT NOT NULL,
    			max_temperature DOUBLE PRECISION NOT NULL,
    			consecutive_breaches_required INTEGER NOT NULL DEFAULT 2,
    			UNIQUE (user_id, city_name)
			);`)

	if err != nil {
		log.Println("error while creating a user_city_thresholds table")
		return err
	}

	_, err = r.db.Exec(`CREATE TABLE IF NOT EXISTS temperature_alerts (
    			id SERIAL,
    			user_id INTEGER REFERENCES users(id),
    			city_name TEXT NOT NULL,
    			temperature DOUBLE PRECISION NOT NULL,
    			threshold DOUBLE PRECISION NOT NULL,
    			consecutive_count INTEGER NOT NULL,
    			alert_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    			email_sent BOOLEAN DEFAULT false
			);

			SELECT create_hypertable('temperature_alerts', 'alert_timestamp', if_not_exists => TRUE);`)

	if err != nil {
		log.Println("error while creating a temperature_alerts table")
		return err
	}

	err = r.Create_threshold_view()

	if err != nil {
		return err
	}

	return nil

}

func (r *AlertRepository) Create_threshold_view() error {
	_, err := r.db.Exec(
		`CREATE MATERIALIZED VIEW IF NOT EXISTS threshold_breach_counts
					WITH (timescaledb.continuous) AS
					SELECT 
    					time_bucket('5 minutes', wr.timestamp) AS bucket,
    					wr.city_name,
    					uct.user_id,
    					COUNT(*) FILTER (WHERE wr.temperature > uct.max_temperature) as breach_count,
    					MAX(wr.temperature) as max_temp,
    					uct.max_temperature as threshold
					FROM weather_records wr
					CROSS JOIN user_city_thresholds uct
				WHERE wr.city_name = uct.city_name
				GROUP BY 
    				time_bucket('5 minutes', wr.timestamp),
    				wr.city_name,
    				uct.user_id,
    				uct.max_temperature;`)

	if err != nil {
		log.Println("error while creating consecutive_threshold_breaches view ")
		return err
	}

	var viewExists bool
	err = r.db.QueryRow(`
	    SELECT EXISTS (
	        SELECT 1 
	        FROM information_schema.views 
	        WHERE table_schema = 'public' 
	          AND table_name = 'threshold_breach_counts'
	    );
	`).Scan(&viewExists)

	if err != nil {
		log.Println("error checking view existence:")
		return err
	}

	if !viewExists {

		_, err = r.db.Exec(`
		SELECT add_continuous_aggregate_policy('threshold_breach_counts',
		start_offset => INTERVAL '1 hour',
		end_offset => INTERVAL '5 minutes',
		schedule_interval => INTERVAL '5 minutes');
		`)

		if err != nil {
			log.Println("error adding aggrgation policy for consecutive_threshold_breaches view ")
			return err
		}
	}

	return nil

}

func (r *AlertRepository) QyeryBreaches_view() (*sql.Rows, error) {

	log.Println("INside query breaches view...")
	// 	rows, err := r.db.Query(`
	// 	WITH new_breaches AS (
	// 		SELECT
	// 			ctb.city_name,
	// 			ctb.user_id,
	// 			ctb.temperature,
	// 			ctb.threshold,
	// 			ctb.consecutive_breaches,
	// 			u.email,
	// 			u.name,
	// 			u.temperature_unit
	// 		FROM consecutive_threshold_breaches ctb
	// 		JOIN users u ON ctb.user_id = u.id
	// 		JOIN user_city_thresholds uct ON
	// 			ctb.user_id = uct.user_id AND
	// 			ctb.city_name = uct.city_name
	// 		WHERE ctb.consecutive_breaches >= uct.consecutive_breaches_required
	// 		AND NOT EXISTS (
	// 			SELECT 1 FROM temperature_alerts ta
	// 			WHERE ta.user_id = ctb.user_id
	// 			AND ta.city_name = ctb.city_name
	// 			AND ta.alert_timestamp > NOW() - INTERVAL '1 hour'
	// 		)
	// 	)
	// 	INSERT INTO temperature_alerts (
	// 		user_id, city_name, temperature, threshold,
	// 		consecutive_count, alert_timestamp
	// 	)
	// 	SELECT
	// 		user_id, city_name, temperature, threshold,
	// 		consecutive_breaches, NOW()
	// 	FROM new_breaches
	// 	RETURNING id, user_id, city_name, temperature, threshold, email
	// `)
	rows, err := r.db.Query(`
		WITH consecutive_breaches AS (
				SELECT 
					bucket,
					city_name,
					user_id,
					breach_count,
					max_temp as temperature,
					threshold,
					LAG(breach_count) OVER (
						PARTITION BY city_name, user_id 
						ORDER BY bucket
					) as prev_breach_count
				FROM threshold_breach_counts
			),
		breach_status AS (
				SELECT 
					bucket,
					city_name,
					user_id,
					temperature,
					threshold,
					CASE 
						WHEN breach_count > 0 AND prev_breach_count > 0 THEN 2
						WHEN breach_count > 0 THEN 1
						ELSE 0
					END as consecutive_breaches
				FROM consecutive_breaches
				WHERE breach_count > 0 OR prev_breach_count > 0
			),
		new_breaches AS (
				SELECT
					bs.city_name,
					bs.user_id,
					bs.temperature,
					bs.threshold,
					bs.consecutive_breaches,
					u.email,
					u.temperature_unit
				FROM breach_status bs
				JOIN users u ON bs.user_id = u.id
				JOIN user_city_thresholds uct ON
					bs.user_id = uct.user_id AND
					bs.city_name = uct.city_name
				WHERE bs.consecutive_breaches >= uct.consecutive_breaches_required
				AND NOT EXISTS (
					SELECT 1 FROM temperature_alerts ta
					WHERE ta.user_id = bs.user_id
					AND ta.city_name = bs.city_name
					AND ta.alert_timestamp > NOW() - INTERVAL '2 minutes'
				)
			),
		inserted_alerts AS (	
			INSERT INTO temperature_alerts (
				user_id, 
				city_name, 
				temperature, 
				threshold,
				consecutive_count, 
				alert_timestamp
			)
			SELECT
				user_id,
				city_name,
				temperature,
				threshold,
				consecutive_breaches,
				NOW()
			FROM new_breaches
			RETURNING id, user_id, city_name, temperature, threshold, consecutive_count, alert_timestamp
		)
			SELECT 
				ia.id,
				ia.user_id,
				ia.city_name,
				ia.temperature,
				ia.threshold,
				nb.email 
			FROM inserted_alerts ia
			JOIN new_breaches nb ON ia.user_id = nb.user_id AND ia.city_name = nb.city_name;	
    
   `)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return rows, nil
}

func (r *AlertRepository) UpdateTemeratureAlert(alertId int) error {

	_, err := r.db.Exec(
		"UPDATE temperature_alerts SET email_sent = true WHERE id = $1",
		alertId,
	)

	return err
}
