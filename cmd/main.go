package main

import (
	routes "example/STRUCTURE/pkg/api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := "8080"

	router := gin.New()
	router.Use(gin.Logger())
	routes.AUTHRoutes(router)
	routes.UserRoutes(router)
	routes.OrganizationRoutes(router)
	router.Run(":" + port)
}
