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

	cfg := config.Load()

	db, err := config.ConnectDB(cfg.DatabaseURL)

	if err != nil {
		log.Println("Failed to connect with database:", err)
	}

	log.Println("Connected to TimescaleDB database")

	//Repo innitialization

	weatherRepo, err := repositories.NewWeatherRepository(db)
	if err != nil {
		log.Println("Failed to initialize new weather repo ", err)
	}

	userRepo := repositories.NewUserRepository(db)

	alertRepo, err := repositories.NewAlertRepository(db)

	if err != nil {
		log.Println("failed to initialize alert repository", err)
	}

	// Initialize OpenWeather client
	weatherClient := openweather.NewClient(cfg.OpenWeatherAPIKey)

	log.Println("Initialized OpenWeather client")

	// Initialize services
	weatherService := services.NewWeatherService(weatherClient, weatherRepo)
	alertService := services.NewAlertService(
		// services.NewAlertService(
		alertRepo,
		"smtp.zoho.com",
		587,
		"tejasthombare@zohomail.in",
		"7iB5#Vu2vXPGamn",
		"tejasthombare@zohomail.in")

	userServices := services.NewUserSerivce(userRepo)

	weatherController := controllers.NewWeatherHandler(weatherService)
	userController := controllers.NewUserController(userServices)

	go weatherService.StartWeatherCollection()
	alertService.StartAlertProcessing()

	routes.WeathterMonitoringRoutes(router, weatherController, userController)

	log.Fatal(router.Run(":" + PORT))
	log.Printf("Starting weather monitoring server on port %s...\n", PORT)

}
