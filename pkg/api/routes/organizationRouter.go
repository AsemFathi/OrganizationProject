package routes

import (
	controller "example/STRUCTURE/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func OrganizationRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/organization", controller.CreateOrganization())
	incomingRoutes.GET("/organization", controller.GetAllOrganizations())
	incomingRoutes.GET("/organization/:org_id", controller.GetOrganization())
	incomingRoutes.PUT("/organization/:org_id", controller.UpdateOrganization())
	incomingRoutes.DELETE("/organization/:org_id", controller.DeleteOrganization())
	incomingRoutes.POST("/organization/:org_id/invite", controller.InviteUserToOrganization())
}
