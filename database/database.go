package database

import (
	"github.com/MonsterYNH/athena/config"
	"gorm.io/gorm"
)

type DatabaseAble interface {
	NewDatabase(*config.DatabaseConfig) (*gorm.DB, error)
}

type Database struct {
	DatabaseAble
	config *config.DatabaseConfig
	db     *gorm.DB
}

func NewDatabaseWithOption(database DatabaseAble, opts ...config.DatabaseConfigOption) (*Database, error) {
	config := &config.DatabaseConfig{}
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
