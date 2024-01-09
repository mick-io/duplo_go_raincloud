package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/mick-io/duplo_go_cloud/internal/api"
	"github.com/mick-io/duplo_go_cloud/internal/database"
	"github.com/mick-io/duplo_go_cloud/internal/handlers"
)

func Initialize(e *echo.Echo, db database.Datastore, client api.WeatherAPIClient) {
	e.GET("/health", handlers.HealthCheckHandler(db))

	e.POST("/locations", handlers.CreateLocation(db, client))
	e.GET("/locations", handlers.ReadLocations(db))
	// e.PUT("/locations/:id", handlers.UpdateLocation(db))
	e.DELETE("/locations/:id", handlers.DeleteLocationByID(db))
	e.DELETE("/locations", handlers.DeleteLocationByLatLong(db))

	e.GET("/forecast", handlers.ReadStoredForecast(db))
	e.PUT("/forecast/latest", handlers.ReadLatestForecast(db, client))
}
