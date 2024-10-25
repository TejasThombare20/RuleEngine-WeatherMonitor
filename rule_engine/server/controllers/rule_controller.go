package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TejasThombare20/rule-engine/models"
	"github.com/TejasThombare20/rule-engine/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		log.Panicln("JSON parsing error", err)
		return
	}

	err := c.service.CreateRule(rule)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error(), "message": "Unable to create rule", "success": false})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "rule created", "success": true})
}

func (c *RuleController) GetRule(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Id is missing", "success": false})
	}

	rule, err := c.service.GetRule(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Rule not found"})
		return
	}

	fmt.Println("rule", rule)
	ctx.JSON(http.StatusOK, gin.H{"rule": rule, "message": "rule fetch succefully", "success": true})
}

func (c *RuleController) GetRules(ctx *gin.Context) {

	rules, err := c.service.GetRules()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch rules"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"rules": rules, "message": "Rules fetched successfully", "success": true})

}

func (c *RuleController) CombineRules(ctx *gin.Context) {
	var input struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description" binding:"required"`
		RuleIDs     []string `json:"rule_ids" binding:"required,min=2"`
	}
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Invalid JSON", "success": true})
		return
	}

	ruleIDs := make([]primitive.ObjectID, len(input.RuleIDs))

	for i, idStr := range input.RuleIDs {
		id, err := primitive.ObjectIDFromHex(idStr)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
			return
		}
		ruleIDs[i] = id
	}

	err := c.service.CombinedRules(input.Name, input.Description, ruleIDs)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to combine rules", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "rules combined successfully", "success": true})
}
