package middleware

import (
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func AdminOnlyMiddleware(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		//технически id не может быть пустым, т.к. перед этим стработало authMiddlware
		id := c.GetString("id")

		UUID, parceErr := uuid.Parse(id)
		if parceErr != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Битый uuid"})
		}

		isAdmin, err := userRepo.CheckUserIsAdmin(UUID)

		//опять-таки, проверки технически не может вернуть факт того, что юзера не существует, потому это сто пудов ошибка бд и тд
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !isAdmin {
			c.AbortWithStatus(http.StatusForbidden)
		} else {
			c.Next()
		}
	}
}
