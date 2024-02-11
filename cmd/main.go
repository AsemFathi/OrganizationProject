package main

import (
	"example/STRUCTURE/pkg/api/middleware"
	routes "example/STRUCTURE/pkg/api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := "8080"

	router := gin.New()
	router.Use(gin.Logger())
	routes.AUTHRoutes(router)

	router.Use(middleware.Authenticate())

	routes.OrganizationRoutes(router)
	routes.UserRoutes(router)
	router.Run(":" + port)
}
