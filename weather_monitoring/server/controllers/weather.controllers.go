package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/TejasThombare20/weather-engine/services"
	"github.com/gin-gonic/gin"
)

type WeatherController struct {
	weatherService *services.WeatherService
}

func NewWeatherHandler(service *services.WeatherService) *WeatherController {
	return &WeatherController{weatherService: service}
}

func (h *WeatherController) GetCityDailySummary(c *gin.Context) {

	cityName := c.Param("city")

	log.Println("city name for daily summary: ", cityName)

	// Bind JSON input to the struct
	if cityName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city name"})
		return
	}

	// today's date for getting daily summary
	todayDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	// todayDate := "2024-10-26"

	summary, err := h.weatherService.GetWeatherSummary(cityName, todayDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// type getDataReqest struct {
// 	City string `json:"city" binding:"optional"`
// }

func (h *WeatherController) GetAllCitydata(c *gin.Context) {
	// var req getDataReqest

	cityName := c.Param("city")

	// if err {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request ", "error": err})
	// 	return
	// }

	// Parse the date from the request

	weatherRecords, error := h.weatherService.GetAllCitydata(cityName)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}

	c.JSON(http.StatusOK, weatherRecords)
}
