package routes

import (
	"time"

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

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	// создание репозиториев (круды для работы с entity)
	userRepo := repository.NewUserRepository(dbConn)
	projectRepo := repository.NewProjectRepository(dbConn)
	slotRepo := repository.NewSlotRepository(dbConn)

	// создание сервисов (бизнес логика)
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	projectService := services.NewProjectService(projectRepo)
	slotService := services.NewSlotService(slotRepo)

	// создание хендлеров
	authHandler := handlers.NewAuthHandler(authService, config)
	userHandler := handlers.NewUserHandler(userService, config)
	skillHandler := handlers.NewSkillsHandler(dbConn)
	projectHandler := handlers.NewProjectHandler(projectService)
	slotHandler := handlers.NewSlotHandler(slotService, projectService)

	api := router.Group("/api")
	{
		// публичные эндпоинты
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/refresh", authHandler.Refresh)
		api.GET("/skills", skillHandler.GetSkills)
		api.GET("/skills/:id", skillHandler.GetSkillByID)

		api.GET("/projects", projectHandler.GetProjects)
		api.GET("/projects/:projectID", projectHandler.GetProjectByID)
		api.GET("/projects/:projectID/slots", slotHandler.GetSlots)

		// защищенные
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(config.JWTSecret, userRepo))
		{
			protected.POST("/projects", projectHandler.CreateProject)
			protected.PUT("/projects/:projectID", projectHandler.UpdateProjectByID)
			protected.DELETE("/projects/:projectID", projectHandler.DeleteProjectByID)

			protected.POST("/projects/:projectID/slots", slotHandler.CreateSlot)
			protected.PUT("/projects/:projectID/slots/:slotID", slotHandler.UpdateSlotByID)

			protected.GET("/users/me", userHandler.Me)
		}

	}

	return router
}
