package main

import (
	"log"

	"github.com/TejasThombare20/rule-engine/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	PORT := "8000"

	router := gin.New()
	router.Use(cors.Default())

	router.Use(gin.Logger())

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	routes.RuleRoutes(router)

	log.Fatal(router.Run(":" + PORT))

}
