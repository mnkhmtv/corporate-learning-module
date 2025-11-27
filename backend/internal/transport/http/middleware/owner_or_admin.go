package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// OwnerOrAdminOnly allows access only to resource owner or admin
func OwnerOrAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Get resource ID from URL parameter
		resourceID := c.Param("id")

		// Allow if user is admin or resource owner
		if role.(string) == "admin" || userID.(string) == resourceID {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		c.Abort()
	}
}
