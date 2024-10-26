package main

import (
	"log"

	openweather "github.com/TejasThombare20/weather-engine/client"
	"github.com/TejasThombare20/weather-engine/config"
	"github.com/TejasThombare20/weather-engine/controllers"
	"github.com/TejasThombare20/weather-engine/routes"

	// "github.com/TejasThombare20/weather-engine/models"
	"github.com/TejasThombare20/weather-engine/repositories"
	"github.com/TejasThombare20/weather-engine/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	PORT := "9000"

	router := gin.New()
	router.Use(cors.Default())

	router.Use(gin.Logger())

	// db, err := config.ConnectDB()
	// if err != nil {
	// 	log.Println("Error connecting to PostgreSQL: ", err)
	// }
	// log.Println("Connected to PostgreSQL database:")
	// db.AutoMigrate(&models.WeatherRecord{})

	cfg := config.Load()
	// if err != nil {
	// 	log.Fatal("Failed to load config:", err)
	// }

	log.Println("Hello")
	// Initialize repository
	repo, err := repositories.NewWeatherRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize repository:", err)
	}

	log.Println("Connected to TimescaleDB database:")

	// Initialize OpenWeather client
	weatherClient := openweather.NewClient(cfg.OpenWeatherAPIKey)

	log.Println("Initialized OpenWeather client")

	// Initialize service
	weatherService := services.NewWeatherService(weatherClient, repo)

	log.Println("Initialized weather service")

	weatherController := controllers.NewWeatherHandler(weatherService)

	// weatherClient := openweather.NewClient("eabed5afa938251da25ac3a124c19c56")
	// weatherRepo := repositories.NewWeatherRepository(db)
	// weatherService := services.NewWeatherService(weatherClient, weatherRepo)

	go weatherService.StartWeatherCollection()

	log.Printf("Starting weather monitoring server on port %s...\n", PORT)

	routes.WeathterMonitoringRoutes(router, weatherController)

	log.Println(router.Run(":" + PORT))

}
