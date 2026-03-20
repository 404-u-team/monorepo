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

	dbConn := db.InitDB(&config)

	db.Migrate(dbConn)

	db.CreateAdmin(dbConn, config)

	router := routes.SetupRoutes(dbConn, &config)

	log.Printf("Server starting on :%s", config.APIPort)
	log.Fatal(http.ListenAndServe(":"+config.APIPort, router))
}
