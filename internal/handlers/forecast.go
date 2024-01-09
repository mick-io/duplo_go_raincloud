package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/mick-io/duplo_go_cloud/internal/api"
	"github.com/mick-io/duplo_go_cloud/internal/database"
	"github.com/mick-io/duplo_go_cloud/internal/models"
)

func ReadStoredForecast(db database.Datastore) echo.HandlerFunc {
	return func(c echo.Context) error {
		locations := []models.LocationRecord{}

		if err := db.Find(&locations); err != nil {
			msg := "Error querying database"
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}

		var wg sync.WaitGroup
		wg.Add(len(locations))
		forecasts := make(chan *models.ReadForecastResponseBody, len(locations))
		errs := make(chan error, len(locations))

		for _, location := range locations {
			go func(location *models.LocationRecord) {
				defer wg.Done()

				forecastRecord := models.ForecastRecord{}
				err := db.Find(&forecastRecord, &models.ForecastRecord{LocationRecordID: location.ID})
				if err != nil {
					errs <- err
					return
				}

				hourlyUnitsRecord := models.HourlyUnitsRecord{}
				err = db.Find(&hourlyUnitsRecord, &models.HourlyUnitsRecord{ForecastRecordID: forecastRecord.ID})
				if err != nil {
					errs <- err
					return
				}

				hourlyRecords := []models.HourlyRecord{}
				err = db.Find(&hourlyRecords, &models.HourlyRecord{ForecastRecordID: forecastRecord.ID})
				if err != nil {
					errs <- err
					return
				}

				forecast := models.ReadForecastResponseBody{
					Latitude:             location.Latitude,
					Longitude:            location.Longitude,
					GenerationtimeMS:     forecastRecord.GenerationtimeMS,
					UTCOffsetSeconds:     forecastRecord.UTCOffsetSeconds,
					Timezone:             forecastRecord.Timezone,
					TimezoneAbbreviation: forecastRecord.TimezoneAbbreviation,
					Elevation:            forecastRecord.Elevation,
					HourlyUnits: models.HourlyUnits{
						Temperature2M: hourlyUnitsRecord.Temperature2MUnit,
						Time:          hourlyUnitsRecord.TimeUnit,
					},
					Hourly: models.Hourly{
						Time:          make([]string, len(hourlyRecords)),
						Temperature2M: make([]float64, len(hourlyRecords)),
					},
				}

				for i, record := range hourlyRecords {
					forecast.Hourly.Time[i] = record.Time
					forecast.Hourly.Temperature2M[i] = record.Temperature2M
				}

				forecasts <- &forecast

			}(&location)
		}

		go func() {
			wg.Wait()
			close(errs)
			close(forecasts)
		}()

		select {
		case err := <-errs:
			if err != nil {
				msg := "Error querying database"
				return echo.NewHTTPError(http.StatusInternalServerError, msg)
			}
		case <-time.After(time.Second * 10):
			msg := "Timeout waiting for goroutines to finish"
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}

		resp := make([]*models.ReadForecastResponseBody, 0)
		for forecast := range forecasts {
			resp = append(resp, forecast)
		}

		return c.JSON(http.StatusOK, &resp)
	}
}

func ReadLatestForecast(db database.Datastore, WeatherAPIClient api.WeatherAPIClient) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: Fetch and upsert latest forecast
		return ReadStoredForecast(db)(c)
	}
}
