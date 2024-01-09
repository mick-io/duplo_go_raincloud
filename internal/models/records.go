package models

import (
	"gorm.io/gorm"
)

type LocationRecord struct {
	gorm.Model
	Latitude        float64
	Longitude       float64
	ForecastRecords []ForecastRecord `gorm:"foreignKey:LocationRecordID"`
}

func NewLocationRecord(loc *Forecast) *LocationRecord {
	return &LocationRecord{
		Latitude:  loc.Latitude,
		Longitude: loc.Longitude,
	}
}

type ForecastRecord struct {
	gorm.Model
	LocationRecordID     uint
	GenerationtimeMS     float64
	UTCOffsetSeconds     int64
	Timezone             string
	TimezoneAbbreviation string
	Elevation            float64
	HourlyRecords        []HourlyRecord    `gorm:"foreignKey:ForecastRecordID"`
	HourlyUnitsRecord    HourlyUnitsRecord `gorm:"foreignKey:ForecastRecordID"`
}

func NewForecastRecords(locationID uint, forecastData *Forecast) *ForecastRecord {
	forecastRecord := &ForecastRecord{
		LocationRecordID:     locationID,
		GenerationtimeMS:     forecastData.GenerationtimeMS,
		UTCOffsetSeconds:     forecastData.UTCOffsetSeconds,
		Timezone:             forecastData.Timezone,
		TimezoneAbbreviation: forecastData.TimezoneAbbreviation,
		Elevation:            forecastData.Elevation,
	}

	return forecastRecord
}

type HourlyRecord struct {
	gorm.Model
	ForecastRecordID uint
	Time             string
	Temperature2M    float64
}

func NewHourlyRecord(forecastRecordID uint, data *Hourly) *[]HourlyRecord {
	var hourlyRecords []HourlyRecord
	for i, time := range data.Time {
		hourlyRecord := HourlyRecord{
			ForecastRecordID: forecastRecordID,
			Time:             time,
			Temperature2M:    data.Temperature2M[i],
		}
		hourlyRecords = append(hourlyRecords, hourlyRecord)
	}
	return &hourlyRecords
}

type HourlyUnitsRecord struct {
	gorm.Model
	ForecastRecordID  uint
	TimeUnit          string
	Temperature2MUnit string
}

func NewHourlyUnitsRecord(forecastRecordID uint, data *HourlyUnits) *HourlyUnitsRecord {
	return &HourlyUnitsRecord{
		ForecastRecordID:  forecastRecordID,
		TimeUnit:          data.Time,
		Temperature2MUnit: data.Temperature2M,
	}
}
