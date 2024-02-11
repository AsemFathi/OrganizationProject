package middleware

import (
	helper "example/STRUCTURE/pkg/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		hasPrefix := strings.HasPrefix(clientToken, "Bearer")

		if clientToken == "" || !hasPrefix {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provides")})
			c.Abort()
			return
		}

		token := strings.SplitAfter(clientToken, " ")[1]

		claims, err := helper.ValidateToken(token)

		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("uid", claims.Uid)
		c.Set("user_type", claims.User_type)
		c.Next()

	}
}
