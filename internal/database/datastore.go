package database

type Datastore interface {
	Find(out interface{}, where ...interface{}) error
	First(out interface{}, where ...interface{}) error
	Create(value interface{}) error
	Save(value interface{}) error
	Delete(value interface{}, where ...interface{}) error
	HealthCheck() error
}
