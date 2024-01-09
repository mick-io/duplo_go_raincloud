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
			Status:   "healthy",
			Database: "up",
		}

		if err := db.HealthCheck(); err != nil {
			status.Status = "unhealthy"
			status.Database = "down"
		}

		if status.Status != "healthy" {
			return c.JSON(http.StatusInternalServerError, status)
		}
		return c.JSON(200, status)
	}
}
