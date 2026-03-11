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
		id, _ := c.Get(UserIdKey)

		isAdmin, err := userRepo.CheckUserIsAdmin(id.(uuid.UUID))

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
