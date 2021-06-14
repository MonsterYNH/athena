package orm

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host         string
	User         string
	Password     string
	DBName       string
	Port         int
	SSLMode      string
	TimeZone     string
	MaxIdleConns int
	MaxOpenConns int
}

type DB struct {
	db    *gorm.DB
	sqlDB *sql.DB
}

var db *DB

func init() {
	database, err := NewDB(&DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "athena",
		Password: "athena",
		DBName:   "athena",
		SSLMode:  "disable",
		TimeZone: "Asia/Shanghai",
	})
	if err != nil {
		panic(err)
	}

	if err := database.db.AutoMigrate(User{}); err != nil {
		panic(err)
	}

	db = database
}

func GetDB() *gorm.DB {
	return db.db
}

func NewDB(config *DBConfig) (*DB, error) {
	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN: fmt.Sprintf(
				"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
				config.Host,
				config.Port,
				config.User,
				config.Password,
				config.DBName,
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

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)

	return &DB{
		db:    db,
		sqlDB: sqlDB,
	}, nil
}

func (db *DB) GetDB() *sql.DB {
	return db.sqlDB
}
