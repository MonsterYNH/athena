package config

import (
	"log"

	"github.com/spf13/viper"
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

type DatabaseConfigOption func(*DatabaseConfig) error

func newDatabaseConfig(v *viper.Viper) (*DatabaseConfig, error) {
	sslMode := v.GetString("database.ssl_mode")
	if len(sslMode) == 0 {
		log.Println("[Info] database.ssl_mode use default value: disable")
		sslMode = "disable"
	}
	timeZone := v.GetString("database.time_zone")
	if len(timeZone) == 0 {
		log.Println("[Info] database.time_zone use default value: Asia/Shanghai")
		timeZone = "Asia/Shanghai"
	}

	return &DatabaseConfig{
		Host:         v.GetString("database.host"),
		Port:         v.GetInt("database.port"),
		User:         v.GetString("database.user"),
		Password:     v.GetString("database.password"),
		Name:         v.GetString("database.name"),
		SSLMode:      sslMode,
		TimeZone:     timeZone,
		MaxIdleConns: v.GetInt("database.max_idle_conns"),
		MaxOpenConns: v.GetInt("database.max_open_conns"),
	}, nil
}
