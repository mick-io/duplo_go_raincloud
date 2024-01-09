package models

type HealthStatusResponseBody struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

type CreateLocationResponseBody struct {
	ID        uint    `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ReadLocationResponseBody struct {
	ID        uint    `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type UpdateLocationResponseBody struct {
	ID        uint    `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type DeleteLocationResponseBody struct {
	ID        uint    `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ReadForecastResponseBody struct {
	Latitude             float64     `json:"latitude" validate:"required"`
	Longitude            float64     `json:"longitude" validate:"required"`
	GenerationtimeMS     float64     `json:"generationtime_ms" validate:"required"`
	UTCOffsetSeconds     int64       `json:"utc_offset_seconds" validate:"required"`
	Timezone             string      `json:"timezone" validate:"required"`
	TimezoneAbbreviation string      `json:"timezone_abbreviation" validate:"required"`
	Elevation            float64     `json:"elevation" validate:"required"`
	HourlyUnits          HourlyUnits `json:"hourly_units" validate:"required"`
	Hourly               Hourly      `json:"hourly" validate:"required"`
}
