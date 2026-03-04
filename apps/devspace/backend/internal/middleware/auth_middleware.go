package middleware

import (
	"net/http"
	"strings"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(JWTSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getAccessToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No authorization token",
			})
			c.Abort()
			return
		}

		userID, err := auth.ValidateJWT([]byte(JWTSecret), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", userID)
		// c.Set("userRole", claims.Role)

		c.Next()
	}
}

func getAccessToken(c *gin.Context) string {
	authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
	if authHeader == "" {
		return ""
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 {
		return ""
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return parts[1]
}
