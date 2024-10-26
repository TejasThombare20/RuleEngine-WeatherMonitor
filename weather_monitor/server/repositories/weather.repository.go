package repositories

import (
	"database/sql"
	"log"

	"github.com/TejasThombare20/weather-engine/models"
	// "gorm.io/gorm"
)

// type WeatherRepository struct {
// 	db *gorm.DB
// }

type WeatherRepository struct {
	db *sql.DB
}

// func NewWeatherRepository(db *gorm.DB) *WeatherRepository {
// 	return &WeatherRepository{db: db}
// }

// func (r *WeatherRepository) SaveWeatherRecord(record *models.WeatherRecord) error {
// 	return r.db.Create(record).Error
// }

func NewWeatherRepository(connStr string) (*WeatherRepository, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	repo := &WeatherRepository{db: db}
	if err := repo.initialize(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *WeatherRepository) initialize() error {
	// Create TimescaleDB extension
	if _, err := r.db.Exec(`CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;`); err != nil {
		log.Println("Error creating TimescaleDB extension:", err)
		return err
	}

	// Create weather_records table and hypertable
	if err := r.createWeatherRecordsTable(); err != nil {
		log.Println("Error creating weather_records table:", err)
		return err
	}

	// Create continuous aggregate
	if err := r.createContinuousAggregate(); err != nil {
		log.Println("Error creating continuous aggregate:", err)
		return err
	}

	return nil
}

func (r *WeatherRepository) createWeatherRecordsTable() error {
	_, err := r.db.Exec(`
        CREATE TABLE IF NOT EXISTS weather_records (
            city_name TEXT NOT NULL,
            temperature DOUBLE PRECISION NOT NULL,
            feels_like DOUBLE PRECISION NOT NULL,
            condition TEXT NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL
        );

        SELECT create_hypertable('weather_records', 'timestamp', 
            if_not_exists => TRUE
        );
    `)
	if err != nil {

		log.Println("error creating weather_records table:", err)
	}
	return err
}

// func (r *WeatherRepository) createContinuousAggregate() error {
// 	_, err := r.db.Exec(`
//         CREATE MATERIALIZED VIEW IF NOT EXISTS weather_daily_summaries_cagg
//         WITH (timescaledb.continuous) AS
//         SELECT
//             time_bucket('1 day', timestamp) AS bucket,
//             city_name,
//             avg(temperature) as avg_temperature,
//             max(temperature) as max_temperature,
//             min(temperature) as min_temperature,
//             mode() WITHIN GROUP (ORDER BY condition) as dominant_condition,
//             jsonb_object_agg(
//                 condition,
//                 count(*)
//             ) as condition_counts,
//             count(*) as total_measurements
//         FROM weather_records
//         GROUP BY time_bucket('1 day', timestamp), city_name;

//         SELECT add_continuous_aggregate_policy('weather_daily_summaries_cagg',
//             start_offset => INTERVAL '2 days',
//             end_offset => INTERVAL '1 hour',
//             schedule_interval => INTERVAL '1 hour');
//     `)
// 	return err
// }

func (r *WeatherRepository) createContinuousAggregate() error {
	_, err := r.db.Exec(`
        CREATE MATERIALIZED VIEW IF NOT EXISTS weather_daily_summaries_cagg
        WITH (timescaledb.continuous) AS
        SELECT
            time_bucket('1 day', timestamp) AS bucket,
            city_name,
            avg(temperature) AS avg_temperature,
            max(temperature) AS max_temperature,
            min(temperature) AS min_temperature,
            mode() WITHIN GROUP (ORDER BY condition) AS dominant_condition,
            count(*) AS total_measurements
        FROM weather_records
        GROUP BY time_bucket('1 day', timestamp), city_name;
    `)

	if err != nil {
		log.Println("error creating a materialized view ", err)
		return err
	}

	// _, err = r.db.Exec(`
	// 		 SELECT remove_continuous_aggregate_policy('weather_daily_summaries_cagg');
	// 		`)

	// // Log any error while removing the existing policy
	// if err != nil {
	// 	log.Println("error removing existing continuous aggregate policy (may not exist): ", err)
	// }

	// Then, add the continuous aggregate policy
	_, err1 := r.db.Exec(`
	    SELECT add_continuous_aggregate_policy('weather_daily_summaries_cagg',
	        start_offset => INTERVAL '3 days',
	        end_offset => INTERVAL '1 hour',
	        schedule_interval => INTERVAL '1 hour');
	`)

	if err1 != nil {
		log.Println("error adding continuous aggregate policy: ", err1)
		return nil
	}

	// 	var exists bool
	// 	err = r.db.QueryRow(`
	//     SELECT EXISTS (
	//         SELECT 1
	//         FROM timescaledb_information.continuous_aggregate_policies
	//         WHERE hypertable_name = 'weather_daily_summaries_cagg'
	//     )
	// `).Scan(&exists)

	// 	if err != nil {
	// 		log.Println("error checking for continuous aggregate policy existence: ", err)
	// 		return err
	// 	}

	// if !exists {
	// 	_, err = r.db.Exec(`
	//     SELECT add_continuous_aggregate_policy('weather_daily_summaries_cagg',
	//         start_offset => INTERVAL '1 month',
	//         end_offset => INTERVAL '1 hour',
	//         schedule_interval => INTERVAL '1 hour');
	// `)
	// 	if err != nil {
	// 		log.Println("error adding continuous aggregate policy: ", err)
	// 		return err
	// 	}
	// }

	// Add the refresh policy in a separate statement
	// _, err = r.db.Exec(`
	//     SELECT add_continuous_aggregate_policy('weather_daily_summaries_cagg',
	//         start_offset => INTERVAL '1 month',
	//         end_offset => INTERVAL '1 h',
	//         schedule_interval => INTERVAL '1 h');
	// `)
	// if err != nil {
	// 	log.Println("error adding continuous aggregate policy: ", err)
	// 	return err
	// }

	return nil
}
func (r *WeatherRepository) SaveWeatherRecord(record *models.WeatherRecord) error {
	_, err := r.db.Exec(`
        INSERT INTO weather_records (city_name, temperature, feels_like, condition, timestamp)
        VALUES ($1, $2, $3, $4, $5)
    `, record.CityName, record.Temperature, record.FeelsLike, record.Condition, record.Timestamp)
	return err
}

// GetDailySummary retrieves the daily summary for a specific city and date
func (r *WeatherRepository) GetDailySummary(cityName string, date string) (*models.DailySummary, error) {
	var summary models.DailySummary
	// var conditionCountsJSON []byte

	log.Println("paramters: ", cityName, date)
	//  WHERE (city_name = $1 OR $1 IS NULL)

	err := r.db.QueryRow(`
        SELECT 
            bucket,
            city_name,
            avg_temperature,
            max_temperature,
            min_temperature,
            dominant_condition,
            total_measurements
        FROM weather_daily_summaries_cagg
		 WHERE city_name = $1
        AND bucket = date_trunc('day', $2::timestamptz)
    `, cityName, date).Scan(
		&summary.Bucket,
		&summary.CityName,
		&summary.AvgTemperature,
		&summary.MaxTemperature,
		&summary.MinTemperature,
		&summary.DominantCondition,
		// &conditionCountsJSON,
		&summary.TotalMeasurements,
	)

	if err != nil {
		log.Println("error retrieving daily summary: ", err)
		return nil, err
	}

	// Parse the JSONB condition counts into map
	// if err := json.Unmarshal(conditionCountsJSON, &summary.ConditionCounts); err != nil {
	// 	log.Println("error unmarshalling condition counts JSON: ", err)
	// 	return nil, err
	// }

	return &summary, nil
}

// get data of all cities or particular city for one day

type WeatherAggregate struct {
	Time      string  `json:"time"`      // For the formatted time
	Delhi     float64 `json:"delhi"`     // Average temperature for Delhi
	Mumbai    float64 `json:"mumbai"`    // Average temperature for Mumbai
	Bengaluru float64 `json:"bengaluru"` // Average temperature for Bengaluru
}

func (r *WeatherRepository) GetCityData(cityName string) (*[]models.WeatherRecord, error) {
	var weather_records []models.WeatherRecord
	var rows *sql.Rows
	var err error

	log.Println("city naeme: ", cityName)

	if cityName == "allcities" {
		// If no cityName is provided, fetch data for all cities
		rows, err = r.db.Query(`
			SELECT
				city_name,
				temperature,
				feels_like,
				condition,
				timestamp
			FROM weather_records
			WHERE timestamp >= NOW() - INTERVAL '1 day'
			ORDER BY timestamp ASC
		`)

	} else {
		// If cityName is provided, filter data by that city
		rows, err = r.db.Query(`
			SELECT
				city_name,
				temperature,
				feels_like,
				condition,
				timestamp
			FROM weather_records
			WHERE timestamp >= NOW() - INTERVAL '1 day'
			AND city_name = $1
		`, cityName)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record models.WeatherRecord

		if err := rows.Scan(&record.CityName, &record.Temperature, &record.FeelsLike, &record.Condition, &record.Timestamp); err != nil {
			return nil, err
		}

		weather_records = append(weather_records, record)
	}
	return &weather_records, nil
}

// models/weather.go

// repository/weather.go
// func (r WeatherRepository) GetCityData(cityName string) (interface{}, error) {
// 	// if cityName == "allcities" {
// 	var records []models.AggregatedWeatherRecord
// 	rows, err := r.db.Query(`
//             WITH hourly_data AS (
//                 SELECT
//                     time_bucket('1 hour', timestamp) AS hourly_time,
//                     city_name,
//                     AVG(temperature) as avg_temp,
//                     AVG(feels_like) as avg_feels_like,
//                     MODE() WITHIN GROUP (ORDER BY condition) as most_frequent_condition
//                 FROM weather_records
//                 WHERE timestamp >= NOW() - INTERVAL '24 hours'
//                 GROUP BY hourly_time, city_name
//             ),
//             pivoted_data AS (
//                 SELECT
//                     hourly_time,
//                     TO_CHAR(hourly_time, 'HH24:MI') as time_str,
//                     MAX(CASE WHEN city_name = 'Delhi' THEN avg_temp END) as delhi_temp,
//                     MAX(CASE WHEN city_name = 'Mumbai' THEN avg_temp END) as mumbai_temp,
//                     MAX(CASE WHEN city_name = 'Bengaluru' THEN avg_temp END) as bengaluru_temp,
//                     MAX(CASE WHEN city_name = 'Delhi' THEN avg_feels_like END) as delhi_feels_like,
//                     MAX(CASE WHEN city_name = 'Mumbai' THEN avg_feels_like END) as mumbai_feels_like,
//                     MAX(CASE WHEN city_name = 'Bengaluru' THEN avg_feels_like END) as bengaluru_feels_like,
//                     MAX(CASE WHEN city_name = 'Delhi' THEN most_frequent_condition END) as delhi_condition,
//                     MAX(CASE WHEN city_name = 'Mumbai' THEN most_frequent_condition END) as mumbai_condition,
//                     MAX(CASE WHEN city_name = 'Bengaluru' THEN most_frequent_condition END) as bengaluru_condition
//                 FROM hourly_data
//                 GROUP BY hourly_time
//             )
//             SELECT
//                 time_str,
//                 ROUND(delhi_temp::numeric, 2) as delhi_temp,
//                 ROUND(mumbai_temp::numeric, 2) as mumbai_temp,
//                 ROUND(bengaluru_temp::numeric, 2) as bengaluru_temp,
//                 ROUND(delhi_feels_like::numeric, 2) as delhi_feels_like,
//                 ROUND(mumbai_feels_like::numeric, 2) as mumbai_feels_like,
//                 ROUND(bengaluru_feels_like::numeric, 2) as bengaluru_feels_like,
//                 delhi_condition,
//                 mumbai_condition,
//                 bengaluru_condition
//             FROM pivoted_data
//             ORDER BY hourly_time DESC;
//         `)
// 	if err != nil {
// 		return nil, fmt.Errorf("error querying aggregated weather data: %v", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var record models.AggregatedWeatherRecord
// 		err := rows.Scan(
// 			&record.TimeStr,
// 			&record.DelhiTemp,
// 			&record.MumbaiTemp,
// 			&record.BengaluruTemp,
// 			&record.DelhiFeelsLike,
// 			&record.MumbaiFeelsLike,
// 			&record.BengaluruFeelsLike,
// 			&record.DelhiCondition,
// 			&record.MumbaiCondition,
// 			&record.BengaluruCondition,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("error scanning aggregated weather record: %v", err)
// 		}
// 		records = append(records, record)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error iterating aggregated weather records: %v", err)
// 	}

// 	return records, nil
// 	// }
// }
