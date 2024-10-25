package controllers

import (
	"net/http"

	"github.com/TejasThombare20/rule-engine/models"
	"github.com/gin-gonic/gin"
)

func (c *RuleController) EvaluateRule(ctx *gin.Context) {

	var req models.EvaluationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Invalid JSON"})
		return
	}

	result, err := c.service.EvaluateRule(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": result})
}
