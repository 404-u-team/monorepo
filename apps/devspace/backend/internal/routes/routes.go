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
	if config.AllowAnyOrigin {
		corsConfig.AllowOriginFunc = func(origin string) bool {
			return true
		}
	} else {
		corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	}

	router.Use(cors.New(corsConfig))

	// создание репозиториев (круды для работы с entity)
	userRepo := repository.NewUserRepository(dbConn)
	projectRepo := repository.NewProjectRepository(dbConn)
	slotRepo := repository.NewSlotRepository(dbConn)
	projectRequestRepo := repository.NewProjectRequestRepository(dbConn)
	ideaRepo := repository.NewIdeaRepository(dbConn)

	// создание сервисов (бизнес логика)
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	projectService := services.NewProjectService(projectRepo)
	slotService := services.NewSlotService(slotRepo, projectRepo)
	projectRequestService := services.NewProjectRequestService(projectRequestRepo, slotRepo, projectRepo)
	ideaService := services.NewIdeaService(ideaRepo)

	// создание хендлеров
	authHandler := handlers.NewAuthHandler(authService, config)
	userHandler := handlers.NewUserHandler(userService)
	skillHandler := handlers.NewSkillsHandler(dbConn)
	projectHandler := handlers.NewProjectHandler(projectService)
	slotHandler := handlers.NewSlotHandler(slotService)
	ideaHandler := handlers.NewIdeaHandler(ideaService, dbConn)
	projectRequestHandler := handlers.NewProjectRequestHandler(projectRequestService)
	testDataHandler := handlers.NewTestDataHandler(services.NewTestDataService(dbConn, config))

	api := router.Group("")
	{
		// публичные эндпоинты
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/refresh", authHandler.Refresh)

		api.GET("/users/:userID", userHandler.GetUserByID)

		api.GET("/skills", skillHandler.GetSkills)
		api.GET("/skills/:id", skillHandler.GetSkillByID)

		api.GET("/projects", projectHandler.GetProjects)
		api.GET("/projects/:projectID", projectHandler.GetProjectByID)
		api.GET("/projects/:projectID/slots", slotHandler.GetSlots)

		api.GET("/ideas", ideaHandler.GetIdeas)
		api.GET("/ideas/:id", ideaHandler.GetIdeaByID)

		// тестовые данные (dev-only)
		api.GET("/generate-test-data", testDataHandler.Start)
		api.GET("/generate-test-data/status", testDataHandler.Status)
		api.GET("/generate-test-data/cancel", testDataHandler.Cancel)

		// защищенные
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(config.JWTSecret, userRepo))
		{
			protected.GET("/users/me", userHandler.GetMe)
			protected.PUT("/users/me", userHandler.UpdateMe)

			protected.POST("/projects", projectHandler.CreateProject)
			protected.PUT("/projects/:projectID", projectHandler.UpdateProjectByID)
			protected.DELETE("/projects/:projectID", projectHandler.DeleteProjectByID)
			protected.GET("/projects/:projectID/requests", projectHandler.GetProjectRequests)
			protected.GET("/users/me/requests", projectHandler.GetUserRequests)

			protected.POST("/projects/:projectID/slots", slotHandler.CreateSlot)
			protected.PUT("/projects/:projectID/slots/:slotID", slotHandler.UpdateSlotByID)
			protected.DELETE("/projects/:projectID/slots/:slotID", slotHandler.DeleteSlotByID)

			protected.PUT("/requests/:projectRequestID/accept", projectRequestHandler.AcceptProjectRequest)
			protected.PUT("/requests/:projectRequestID/reject", projectRequestHandler.RejectProjectRequest)

			protected.PUT("/ideas/:ideaID", ideaHandler.UpdateIdeaByID)

			protected.POST("/projects/:projectID/slots/:slotID/apply", projectRequestHandler.CreateProjectRequestApply)
			protected.POST("/projects/:projectID/slots/:slotID/invite", projectRequestHandler.CreateProjectRequestInvite)

			protected.POST("/users/me/skills", skillHandler.AddSkillToSelf)
			protected.DELETE("/users/me/skills/:id", skillHandler.DeleteSelfSkill)

			protected.POST("/ideas", ideaHandler.AddIdea)
			protected.DELETE("/ideas/:id", ideaHandler.DeleteIdeaByID)
		}

		//только для админов
		adminOnly := api.Group("")
		adminOnly.Use(middleware.AuthMiddleware(config.JWTSecret, userRepo), middleware.AdminOnlyMiddleware(userRepo))
		{
			adminOnly.POST("/skills", skillHandler.CreateSkill)
			adminOnly.DELETE("/skills/:id", skillHandler.DeleteSkill)
		}

	}

	return router
}
