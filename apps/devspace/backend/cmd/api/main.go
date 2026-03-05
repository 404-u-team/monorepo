package main

import (
	"log"
	"net/http"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/db"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/routes"
)

func main() {
	config := config.LoadConfig()
	logger := log.Logger{}

	dbConn := db.InitDB(&config)

	db.Migrate(dbConn)

	router := routes.SetupRoutes(dbConn, &config, &logger)

	log.Printf("Server starting on :%s", config.APIPort)
	log.Fatal(http.ListenAndServe(":"+config.APIPort, router))
}
