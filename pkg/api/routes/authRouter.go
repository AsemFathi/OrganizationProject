package routes

import (
	controller "example/STRUCTURE/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func AUTHRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
}
