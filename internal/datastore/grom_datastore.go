package datastore

import (
	"gorm.io/gorm"

	"github.com/mick-io/duplo_go_cloud/internal/database"
)

// GormDatastore is an implementation of the Datastore interface using GORM.
type GormDatastore struct {
	db *gorm.DB
}

// NewGormDatastore creates a new GormDatastore with the given *gorm.DB instance.
func NewGormDatastore(db *gorm.DB) database.Datastore {
	return &GormDatastore{db: db}
}

// Find retrieves records that match the given conditions and stores them in 'out'.
// Examples:
//
// Using a map:
//
//	users := []User{}
//	ds.Find(&users, map[string]interface{}{"name": "mick"})
//
// Using a struct:
//
//	users := []User{}
//	ds.Find(&users, &User{Name: "mick"})
//
// Using a string:
//
//	users := []User{}
//	ds.Find(&users, "name = ?", "mick")
//
// Using a slice:
//
//	users := []User{}
//	ds.Find(&users, "name IN (?)", []string{"mick", "mick 2"})
func (g *GormDatastore) Find(out interface{}, where ...interface{}) error {
	result := g.db.Find(out, where...)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// First retrieves the first record that matches the given conditions and stores it in 'out'.
// Examples:
//
// Using a map:
//
//	user := User{}
//	ds.First(&user, map[string]interface{}{"name": "mick"})
//
// Using a struct:
//
//	user := User{}
//	ds.First(&user, &User{Name: "mick"})
//
// Using a string:
//
//	user := User{}
//	ds.First(&user, "name = ?", "mick")
func (g *GormDatastore) First(out interface{}, where ...interface{}) error {
	result := g.db.First(out, where...)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Create inserts the given value into the database.
// For example:
//
//	user := User{Name: "mick"}
//	ds.Create(&user)
//
// This will create a new user with the name 'mick'.
func (g *GormDatastore) Create(value interface{}) error {
	result := g.db.Create(value)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Save updates the given record in the database. If the record does not exist, it creates a new one.
// For example:
//
//	user := User{Name: "mick"}
//	ds.First(&user)
//	user.Name = "not mick"
//	ds.Save(&user)
//
// This will find the first user with the name 'mick', change its name to 'not mick', and save the change.
func (g *GormDatastore) Save(value interface{}) error {
	result := g.db.Save(value)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete removes records that match the given conditions from the database.
// If an error occurs during the operation, it returns the error.
// Examples:
//
// Using a map:
//
//	ds.Delete(&User{}, map[string]interface{}{"name": "mick"})
//
// Using a struct:
//
//	ds.Delete(&User{}, &User{Name: "mick"})
//
// Using a string:
//
//	ds.Delete(&User{}, "name = ?", "mick")
func (g *GormDatastore) Delete(value interface{}, where ...interface{}) error {
	result := g.db.Delete(value, where...)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// HealthCheck checks the health of the database connection.
func (g *GormDatastore) HealthCheck() error {
	sqlDB, err := g.db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	return nil
}
