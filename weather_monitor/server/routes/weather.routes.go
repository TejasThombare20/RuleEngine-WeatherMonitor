package routes

import (
	"github.com/TejasThombare20/weather-engine/controllers"
	"github.com/gin-gonic/gin"
)

func WeathterMonitoringRoutes(Routes *gin.Engine, weatherController *controllers.WeatherController) {

	Routes.GET("/getsummary/:city", weatherController.GetCityDailySummary)
	Routes.GET("/getRecords/:city", weatherController.GetAllCitydata)
	// Routes.GET("/rules/:id", ruleController.GetRule)
	// Routes.POST("/rules/combine", ruleController.CombineRules)
	// Routes.GET("/rules", ruleController.GetRules)
	// Routes.POST("/evaluate", ruleController.EvaluateRule)

}
