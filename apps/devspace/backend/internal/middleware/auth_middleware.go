package middleware

import (
	"net/http"
	"strings"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/auth"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(JWTSecret string, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getAccessToken(c)
		if token == "" {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		userID, err := auth.ValidateJWT([]byte(JWTSecret), token)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		exists, err := userRepo.IsUserExistByID(userID)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}
		if !exists {
			c.Status(http.StatusUnauthorized)
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
