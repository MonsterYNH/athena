package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDataBase struct{}

func (database *PostgresDataBase) NewDatabase(config *DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN: fmt.Sprintf(
				"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
				config.Host,
				config.Port,
				config.User,
				config.Password,
				config.Name,
				config.SSLMode,
				config.TimeZone,
			),
			PreferSimpleProtocol: true,
		}),
	)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if config.MaxIdleConns != 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	} else {
		log.Println("db max idle conns using default 100")
		sqlDB.SetMaxIdleConns(100)
	}

	if config.MaxOpenConns != 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	} else {
		log.Println("db max open conns using default 100")
		sqlDB.SetMaxOpenConns(100)
	}
	return db, nil
}
