package controllers

import (
	"net/http"

	"github.com/TejasThombare20/weather-engine/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserSerivce
}

func NewUserController(service *services.UserSerivce) *UserController {
	return &UserController{userService: service}
}

type AlertDataRequest struct {
	Email             string             `json:"email"`
	ConsecutiveAlerts int                `json:"consecutiveAlerts"`
	Temperatures      map[string]float64 `json:"Thres_temperatue"` // Use float64 if needed
}

func (h *UserController) CreateUser(c *gin.Context) {

	var req AlertDataRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.AddUserwithCityThreashold(req.Email, "celcius", req.Temperatures, req.ConsecutiveAlerts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "err": err})

	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})

}
