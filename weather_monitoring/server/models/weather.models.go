package models

import (
	"time"
	// "gorm.io/gorm"
)

type WeatherRecord struct {
	CityName    string    `json:"city_name"`
	Temperature float64   `json:"temperature"`
	FeelsLike   float64   `json:"feels_like"`
	Condition   string    `json:"condition"`
	Timestamp   time.Time `json:"timestamp"`
}

type DailySummary struct {
	Bucket            time.Time      `json:"date"`
	CityName          string         `json:"city_name"`
	AvgTemperature    float64        `json:"avg_temperature"`
	MaxTemperature    float64        `json:"max_temperature"`
	MinTemperature    float64        `json:"min_temperature"`
	DominantCondition string         `json:"dominant_condition"`
	ConditionCounts   map[string]int `json:"condition_counts"`
	TotalMeasurements int            `json:"total_measurements"`
}

type AggregatedWeatherRecord struct {
	TimeStr            string  `json:"time"`
	DelhiTemp          float64 `json:"delhi_temp"`
	MumbaiTemp         float64 `json:"mumbai_temp"`
	BengaluruTemp      float64 `json:"bengaluru_temp"`
	DelhiFeelsLike     float64 `json:"delhi_feels_like"`
	MumbaiFeelsLike    float64 `json:"mumbai_feels_like"`
	BengaluruFeelsLike float64 `json:"bengaluru_feels_like"`
	DelhiCondition     string  `json:"delhi_condition"`
	MumbaiCondition    string  `json:"mumbai_condition"`
	BengaluruCondition string  `json:"bengaluru_condition"`
}
