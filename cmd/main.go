package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/mick-io/duplo_go_cloud/internal/api"
	"github.com/mick-io/duplo_go_cloud/internal/config"
	"github.com/mick-io/duplo_go_cloud/internal/database"
	"github.com/mick-io/duplo_go_cloud/internal/datastore"
	"github.com/mick-io/duplo_go_cloud/internal/routes"
)

func main() {
	cfgFP := flag.String("config", "./config/dev.toml", "Path to the config file")
	flag.Parse()

	cfg, err := config.Load(*cfgFP)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := database.Initialize(cfg.Database)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	client := api.NewClient(cfg.API.ForecastAPIBaseURL)
	store := datastore.NewGormDatastore(db)
	e := echo.New()

	routes.Initialize(e, store, client)
	e.Start(":" + strconv.Itoa(cfg.Server.Port))
}
