package routes

import (
	controller "example/STRUCTURE/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func AUTHRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/signup", controller.Signup())
	incomingRoutes.POST("/signin", controller.Signin())
}
