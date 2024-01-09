package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/mick-io/duplo_go_cloud/internal/api"
	"github.com/mick-io/duplo_go_cloud/internal/database"
	"github.com/mick-io/duplo_go_cloud/internal/models"
)

func CreateLocation(db database.Datastore, WeatherAPIClient api.WeatherAPIClient) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Location data validation
		var body models.CreateLocationRequestBody
		if err := c.Bind(&body); err != nil {
			msg := fmt.Sprintf("Failed to parse request body: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, msg)
		}

		if err := body.Validate(); err != nil {
			msg := fmt.Sprintf("Failed to validate location data: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, msg)
		}

		// Checking for conflicting location
		record := models.LocationRecord{}
		if err := db.Find(&record, &models.LocationRecord{Latitude: body.Latitude, Longitude: body.Longitude}); err != nil {
			msg := fmt.Sprintf("Error querying database: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}
		if record.ID != 0 {
			return echo.NewHTTPError(http.StatusConflict, models.CreateLocationResponseBody{
				ID:        record.ID,
				Latitude:  record.Latitude,
				Longitude: record.Longitude,
			})
		}

		// Storing location
		loc := &models.LocationRecord{
			Latitude:  body.Latitude,
			Longitude: body.Longitude,
		}
		if err := db.Create(&loc); err != nil {
			msg := fmt.Sprintf("Error storing location: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}

		// Fetching forecast data
		resp := models.Forecast{}
		opts := api.ForecastOptions{
			Latitude:  strconv.FormatFloat(body.Latitude, 'f', 6, 64),
			Longitude: strconv.FormatFloat(body.Longitude, 'f', 6, 64),
		}
		if err := WeatherAPIClient.GetForecast(opts, &resp); err != nil {
			msg := fmt.Sprintf("Error getting forecast: %v", err)
			return echo.NewHTTPError(http.StatusBadGateway, msg)
		}

		// Validating forecast data
		if err := models.ValidateForecast(resp); err != nil {
			msg := fmt.Sprintf("Error validating forecast: %v", err)
			return echo.NewHTTPError(http.StatusBadGateway, msg)
		}

		// Storing forecast data
		forecast := models.NewForecastRecords(loc.ID, &resp)
		if err := db.Save(forecast); err != nil {
			msg := fmt.Sprintf("Error storing forecast: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}
		hourly := models.NewHourlyRecord(forecast.Model.ID, &resp.Hourly)
		if err := db.Create(&hourly); err != nil {
			msg := fmt.Sprintf("Error storing hourly records: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}
		units := models.NewHourlyUnitsRecord(forecast.Model.ID, &resp.HourlyUnits)
		if err := db.Save(units); err != nil {
			msg := fmt.Sprintf("Error storing unit records: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}

		// Responding with location data
		return c.JSON(http.StatusOK, models.CreateLocationResponseBody{
			ID:        loc.ID,
			Latitude:  loc.Latitude,
			Longitude: loc.Longitude,
		})
	}
}

func ReadLocations(db database.Datastore) echo.HandlerFunc {
	return func(c echo.Context) error {
		records := []models.LocationRecord{}
		if err := db.Find(&records); err != nil {
			msg := fmt.Sprintf("Error querying database: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}

		resp := make([]models.ReadLocationResponseBody, len(records))
		for i, record := range records {
			resp[i] = models.ReadLocationResponseBody{
				ID:        record.ID,
				Latitude:  record.Latitude,
				Longitude: record.Longitude,
			}
		}

		return c.JSON(http.StatusOK, resp)
	}
}

// func UpdateLocation(db database.Datastore) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var body models.UpdateLocationRequestBody
// 		if err := c.Bind(&body); err != nil {
// 			msg := fmt.Sprintf("Failed to parse request body: %v", err)
// 			return echo.NewHTTPError(http.StatusBadRequest, msg)
// 		}
// 		if err := body.Validate(); err != nil {
// 			msg := fmt.Sprintf("Failed to validate location data: %v", err)
// 			return echo.NewHTTPError(http.StatusBadRequest, msg)
// 		}

// 		idParam := c.Param("id")
// 		id, err := strconv.Atoi(idParam)
// 		if err != nil {
// 			msg := fmt.Sprintf("Invalid id parameter: %v", err)
// 			return echo.NewHTTPError(http.StatusBadRequest, msg)
// 		}

// 		var record models.LocationRecord
// 		if err := db.Find(&record, id); err != nil {
// 			msg := fmt.Sprintf("Error querying database: %v", err)
// 			return echo.NewHTTPError(http.StatusInternalServerError, msg)
// 		}

// 		if body.Latitude != nil {
// 			record.Latitude = *body.Latitude
// 		}
// 		if body.Longitude != nil {
// 			record.Longitude = *body.Longitude
// 		}

// 		if err := db.Save(&record); err != nil {
// 			msg := fmt.Sprintf("Error updating location: %v", err)
// 			return echo.NewHTTPError(http.StatusInternalServerError, msg)
// 		}

// 		return c.JSON(http.StatusOK, models.UpdateLocationResponseBody{
// 			ID:        record.ID,
// 			Latitude:  record.Latitude,
// 			Longitude: record.Longitude,
// 		})
// 	}
// }

func DeleteLocationByID(db database.Datastore) echo.HandlerFunc {
	return func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			msg := fmt.Sprintf("Invalid id parameter: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, msg)
		}

		var record models.LocationRecord
		if err := db.Find(&record, id); err != nil {
			msg := fmt.Sprintf("Error querying database: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}

		if record.ID == 0 {
			msg := fmt.Sprintf("Location not found w/ID: %v", id)
			return echo.NewHTTPError(http.StatusNotFound, msg)
		}

		if err := db.Delete(&record); err != nil {
			msg := fmt.Sprintf("Error deleting location: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}

		return c.JSON(http.StatusNoContent, models.DeleteLocationResponseBody{
			ID:        record.ID,
			Latitude:  record.Latitude,
			Longitude: record.Longitude,
		})
	}
}

func DeleteLocationByLatLong(db database.Datastore) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Validating input
		latitudeParam := c.QueryParam("latitude")
		longitudeParam := c.QueryParam("longitude")

		latitude, err := strconv.ParseFloat(latitudeParam, 64)
		if err != nil {
			msg := fmt.Sprintf("Invalid latitude parameter: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, msg)
		}

		longitude, err := strconv.ParseFloat(longitudeParam, 64)
		if err != nil {
			msg := fmt.Sprintf("Invalid longitude parameter: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, msg)
		}

		// Ensuring record exist
		record := models.LocationRecord{}
		if err := db.Find(&record, &models.LocationRecord{Latitude: latitude, Longitude: longitude}); err != nil {
			msg := fmt.Sprintf("Error querying database: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}
		if record.Model.ID == 0 {
			msg := fmt.Sprintf("Location not found w/latitude: %v and longitude: %v", latitude, longitude)
			return echo.NewHTTPError(http.StatusNotFound, msg)
		}

		// Deleting record
		if err := db.Delete(&record, "latitude = ? AND longitude = ?", latitude, longitude); err != nil {
			msg := fmt.Sprintf("Error deleting locations: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}

		return c.JSON(http.StatusNoContent, nil)
	}
}
