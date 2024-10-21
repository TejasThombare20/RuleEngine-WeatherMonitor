package main

import (
	"log"

	"github.com/TejasThombare20/rule-engine/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	PORT := "8000"

	router := gin.New()
	router.Use(cors.Default())

	router.Use(gin.Logger())

	routes.RuleRoutes(router)

	log.Fatal(router.Run(":" + PORT))

}
