package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matiuskm/go-commerce/db"
	"github.com/matiuskm/go-commerce/models"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDAny, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		userID := userIDAny.(uint)

		var user models.User
		if err := db.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
			c.Abort()
			return
		}

		c.Next()
	}
}