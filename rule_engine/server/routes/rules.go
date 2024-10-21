package routes

import (
	"github.com/TejasThombare20/rule-engine/controllers"
	"github.com/TejasThombare20/rule-engine/repositories"
	"github.com/TejasThombare20/rule-engine/services"
	"github.com/gin-gonic/gin"
)

func RuleRoutes(Routes *gin.Engine) {

	ruleRepo := repositories.NewRuleRepository()
	ruleService := services.NewRuleService(ruleRepo)
	ruleController := controllers.NewRuleController(ruleService)

	Routes.POST("/rules", ruleController.CreateRule)
	Routes.GET("/rules/:id", ruleController.GetRule)
	Routes.POST("/rules/combine", ruleController.CombineRules)

}
