package services

import (
	"log"
	"time"

	openweather "github.com/TejasThombare20/weather-engine/client"
	"github.com/TejasThombare20/weather-engine/models"
	"github.com/TejasThombare20/weather-engine/repositories"
)

type WeatherService struct {
	weatherClient *openweather.Client
	weatherRepo   *repositories.WeatherRepository
	cities        []string
}

func NewWeatherService(client *openweather.Client, repo *repositories.WeatherRepository) *WeatherService {
	return &WeatherService{
		weatherClient: client,
		weatherRepo:   repo,
		cities: []string{
			"Delhi",
			"Mumbai",
			"Bengaluru",
		},
	}
}

func (s *WeatherService) StartWeatherCollection() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	log.Println("Starting weather collection...")

	for {
		s.collectWeatherData()
		<-ticker.C
	}
}

func (s *WeatherService) collectWeatherData() {
	log.Println("entering weather collection loop...")
	for _, city := range s.cities {
		weather, err := s.weatherClient.GetWeather(city)
		if err != nil {
			log.Printf("Error fetching weather for %s: %v", city, err)
			continue
		}

		var condition string
		if len(weather.Weather) > 0 {
			condition = weather.Weather[0].Main
		} else {
			condition = "Unknown" // Default or handle accordingly
		}

		record := &models.WeatherRecord{
			CityName:    city,
			Temperature: weather.Main.Temp - 273.15, // Convert Kelvin to Celsius
			FeelsLike:   weather.Main.FeelsLike - 273.15,
			Condition:   condition,
			Timestamp:   time.Now(),
		}

		if err := s.weatherRepo.SaveWeatherRecord(record); err != nil {
			log.Printf("Error saving weather record for %s: %v", city, err)
		}
	}
}

// GetWeatherSummary returns daily summary for a specific city and date
func (s *WeatherService) GetWeatherSummary(cityName string, date string) (*models.DailySummary, error) {
	return s.weatherRepo.GetDailySummary(cityName, date)
}

// getall  data of one daya as well as get all data of one day for one city
func (s *WeatherService) GetAllCitydata(cityName string) (*[]models.WeatherRecord, error) {
	return s.weatherRepo.GetCityData(cityName)
}
