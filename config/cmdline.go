package config

import (
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	// config
	pflag.String("config.file", "", "config with yaml file")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

type Config struct {
	*DatabaseConfig
	*ServiceConfig
	*RegistryConfig
}

func GetConfig() (*Config, error) {
	configFile := viper.GetString("config.file")
	log.Println("config file", configFile)
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}

	databaseConfig, err := newDatabaseConfig(viper.GetViper())
	if err != nil {
		return nil, err
	}
	config.DatabaseConfig = databaseConfig

	serverConfig, err := newServiceConfig(viper.GetViper())
	if err != nil {
		return nil, err
	}
	config.ServiceConfig = serverConfig

	registryConfig, err := newRegistryConfig(viper.GetViper())
	if err != nil {
		return nil, err
	}
	config.RegistryConfig = registryConfig

	return config, nil
}
