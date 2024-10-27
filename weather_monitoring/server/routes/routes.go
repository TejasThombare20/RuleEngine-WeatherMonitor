package routes

import (
	"github.com/TejasThombare20/weather-engine/controllers"
	"github.com/gin-gonic/gin"
)

func WeathterMonitoringRoutes(Routes *gin.Engine, weatherController *controllers.WeatherController, userController *controllers.UserController) {

	Routes.GET("/getsummary/:city", weatherController.GetCityDailySummary)
	Routes.GET("/getRecords/:city", weatherController.GetAllCitydata)
	Routes.POST("/createuser", userController.CreateUser)

}
