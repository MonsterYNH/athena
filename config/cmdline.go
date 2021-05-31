package config

import (
	"fmt"
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	// config
	pflag.String("config.file", "", "config with yaml file")
	pflag.String("config.registry", "", "config with registry")

	// server
	pflag.String("host", "0.0.0.0", "server host")
	pflag.Int("port", 8080, "server port")
	pflag.Bool("enablehttp", true, "enable http gateway")

	// database
	pflag.String("db", "", "database uri eg: postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]")

	// logger
	pflag.String("logger", "", "logger uri")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

func readConfig() error {
	fileConfig := viper.GetString("config.file")
	registryConfig := viper.GetString("config.registry")

	if len(fileConfig) > 0 && len(registryConfig) > 0 {
		return fmt.Errorf("ERROR: file config and registry config both exist, except only one")
	}

	if len(fileConfig) > 0 {
		log.Println(fmt.Sprintf("INFO: use file config: %s", fileConfig))
		// TODO read file config
	} else {
		log.Println(fmt.Sprintf("INFO: use registry config: %s", registryConfig))
		// TODO read registry config
	}

	return nil
}

func readConfigWithFile(filePath string) (Config, error) {
	// TODO 可支持任意配置文件格式
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	return Config{}, nil
}
