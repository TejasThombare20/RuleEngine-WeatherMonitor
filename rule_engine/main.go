package main

import (
	"log"

	"github.com/TejasThombare20/rule-engine/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	PORT := "8085"

	router := gin.New()

	router.Use(gin.Logger())

	routes.RuleRoutes(router)

	log.Fatal(router.Run(":" + PORT))

}
