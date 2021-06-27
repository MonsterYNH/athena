package database

import (
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	SSLMode      string
	TimeZone     string
	MaxIdleConns int
	MaxOpenConns int
}

type DatabaseAble interface {
	NewDatabase(*DatabaseConfig) (*gorm.DB, error)
}

type DatabaseConfigOption func(*DatabaseConfig) error

type Database struct {
	DatabaseAble
	config *DatabaseConfig
	db     *gorm.DB
}

func NewDatabaseWithOption(database DatabaseAble, opts ...DatabaseConfigOption) (*Database, error) {
	config := &DatabaseConfig{}
	for _, option := range opts {
		if err := option(config); err != nil {
			return nil, err
		}
	}

	db, err := database.NewDatabase(config)
	if err != nil {
		return nil, err
	}

	return &Database{
		DatabaseAble: database,
		config:       config,
		db:           db,
	}, nil
}

func (database *Database) GetConnect() *gorm.DB {
	return database.db
}
