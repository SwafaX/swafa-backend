package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/calvinnle/todo-app/initializers"
	"github.com/calvinnle/todo-app/models"
	"github.com/calvinnle/todo-app/utils"
	"github.com/gin-gonic/gin"
)

func DeserializeUSer() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.Request.Header.Get("Authorization")

		if authorizationHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "fail",
				"error":  "authorization header is required",
			})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 || fields[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "fail",
				"error":  "invalid authorization header format, expected 'Bearer <token>'",
			})
			return
		}

		access_token := fields[1]

		// cant reach to this issue
		// if access_token == "" {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		// 		"status": "fail",
		// 		"error": "access token is required",
		// 	})
		// 	return
		// }

		config, _ := initializers.LoadConfig(".")

		sub, err := utils.ValidateToken(access_token, config.AccessTokenPublic)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "fail",
				"message": "Invalid or expired token",
			})
			return
		}

		var user models.User
		result := initializers.DB.Select("id").First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  "fail",
				"message": "User not found or no longer exists",
			})
			return
		}

		c.Set("currentUser", user)
		c.Next()
	}
}
