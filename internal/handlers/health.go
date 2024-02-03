package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/mick-io/duplo_go_cloud/internal/database"
	"github.com/mick-io/duplo_go_cloud/internal/models"
)

func HealthCheckHandler(db database.Datastore) echo.HandlerFunc {
	return func(c echo.Context) error {
		status := models.HealthStatusResponseBody{
			Status:     "OK",
			Database:   "OK",
			WeatherAPI: "OK",
		}

		if err := db.HealthCheck(); err != nil {
			status.Database = "ERROR"
		}

		if status.Status != "ERROR" {
			return c.JSON(http.StatusServiceUnavailable, status)
		}
		return c.JSON(200, status)
	}
}
