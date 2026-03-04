package routes

import (
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/handlers"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/middleware"
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
	projectRepo := repository.NewProjectRepository(dbConn)

	// создание сервисов (бизнес логика)
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	projectService := services.NewProjectService(projectRepo)

	// создание хендлеров
	authHandler := handlers.NewAuthHandler(authService, config)
	userHandler := handlers.NewUserHandler(userService, config)
	projectHandler := handlers.NewProjectHandler(projectService, config)

	api := router.Group("/api")
	{
		// публичные эндпоинты
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/refresh", authHandler.Refresh)

		// защищенные
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(config.JWTSecret, userRepo))
		{
			protected.POST("/projects", projectHandler.CreateProject)

			protected.GET("/users/me", userHandler.Me)
		}

	}

	// router.POST("/api/users/create", rest.RegisterUser)

	return router
}
