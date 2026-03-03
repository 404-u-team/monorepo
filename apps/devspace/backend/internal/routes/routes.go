package routes

import (
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/handlers"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/repository"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(dbConn *gorm.DB, config *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(cors.Default())

	// создание репозиториев (круды для работы с entity)
	userRepo := repository.NewUserRepository(dbConn)

	// создание сервисов (бизнес логика)
	authService := services.NewAuthService(userRepo)

	// создание хендлеров
	authHandler := handlers.NewAuthHandler(authService, config)

	api := router.Group("/api")
	{
		api.POST("/register", authHandler.Register)
	}

	// router.POST("/api/users/create", rest.RegisterUser)

	return router
}
