package controllers

import (
	"net/http"

	"github.com/TejasThombare20/rule-engine/models"
	"github.com/TejasThombare20/rule-engine/services"
	"github.com/gin-gonic/gin"
)

type RuleController struct {
	service *services.RuleService
}

func NewRuleController(service *services.RuleService) *RuleController {
	return &RuleController{service: service}
}

func (c *RuleController) CreateRule(ctx *gin.Context) {

	var rule models.Rule

	if err := ctx.BindJSON(&rule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error(), "message": "JSON parsing error", "success": false})
		return
	}

	err := c.service.CreateRule(rule)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error(), "message": "Unable to create rule", "success": false})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "rule created", "success": true})
}
