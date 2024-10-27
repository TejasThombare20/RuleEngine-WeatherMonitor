package models

import "time"

type TemperatureUnit string

const (
	Celsius    TemperatureUnit = "celsius"
	Fahrenheit TemperatureUnit = "fahrenheit"
)

type User struct {
	ID              int             `json:"id"`
	Email           string          `json:"email"`
	TemperatureUnit TemperatureUnit `json:"temperature_unit"`
}

type CityThreshold struct {
	ID                          int     `json:"id"`
	UserID                      int     `json:"user_id"`
	CityName                    string  `json:"city_name"`
	MaxTemperature              float64 `json:"max_temperature"`
	ConsecutiveBreachesRequired int     `json:"consecutive_breaches_required"`
}

type TemperatureAlert struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	CityName         string    `json:"city_name"`
	Temperature      float64   `json:"temperature"`
	Threshold        float64   `json:"threshold"`
	ConsecutiveCount int       `json:"consecutive_count"`
	AlertTimestamp   time.Time `json:"alert_timestamp"`
	EmailSent        bool      `json:"email_sent"`
}
