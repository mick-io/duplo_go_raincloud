package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mick-io/duplo_go_cloud/internal/config"
	"github.com/mick-io/duplo_go_cloud/internal/models"
)

var DB *gorm.DB

func Initialize(dbCfg *config.DatabaseConfig) (*gorm.DB, error) {
	var err error
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.Name)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB.AutoMigrate(&models.LocationRecord{})
	DB.AutoMigrate(&models.ForecastRecord{})
	DB.AutoMigrate(&models.HourlyRecord{})
	DB.AutoMigrate(&models.HourlyUnitsRecord{})

	return DB, nil
}
