package models

import "github.com/go-playground/validator"

func ValidateForecast(f Forecast) error {
	validate := validator.New()
	validate.RegisterStructValidation(func(sl validator.StructLevel) {
		hourly := sl.Current().Interface().(Hourly)
		if len(hourly.Time) != len(hourly.Temperature2M) {
			sl.ReportError(hourly.Time, "Time", "time", "len", "")
			sl.ReportError(hourly.Temperature2M, "Temperature2M", "temperature_2m", "len", "")
		}
	}, Hourly{})

	err := validate.Struct(f)
	if err != nil {
		return err
	}
	return nil
}

type Forecast struct {
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

type Hourly struct {
	Time          []string  `json:"time" validate:"required"`
	Temperature2M []float64 `json:"temperature_2m" validate:"required"`
}

type HourlyUnits struct {
	Time          string `json:"time" validate:"required"`
	Temperature2M string `json:"temperature_2m" validate:"required"`
}
