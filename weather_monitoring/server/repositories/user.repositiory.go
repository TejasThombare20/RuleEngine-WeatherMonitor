package repositories

import (
	"database/sql"

	"github.com/TejasThombare20/weather-engine/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(email string, temperature_unit string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(`
		INSERT INTO users (email, temperature_unit)
		VALUES ($1, $2)
		RETURNING id, email, temperature_unit`, email, temperature_unit).Scan(
		&user.ID, &user.Email, &user.TemperatureUnit,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) SetCityThreshold(userID int, cityName string, maxTemp float64, consecutiveAlerts int) error {
	_, err := r.db.Exec(`
        INSERT INTO user_city_thresholds 
        (user_id, city_name, max_temperature, consecutive_breaches_required)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id, city_name)
        DO UPDATE SET 
            max_temperature = EXCLUDED.max_temperature,
            consecutive_breaches_required = EXCLUDED.consecutive_breaches_required
    `, userID, cityName, maxTemp, consecutiveAlerts)
	return err
}
