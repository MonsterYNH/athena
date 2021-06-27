package config

import (
	"errors"
	"fmt"
	"net"

	"github.com/spf13/viper"
)

type ServiceConfig struct {
	ServiceName string `json:"service_name"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	EnableHTTP  bool   `json:"enable_http"`
}

func newServiceConfig(v *viper.Viper) (*ServiceConfig, error) {
	serviceName := v.GetString("service.service_name")
	if len(serviceName) == 0 {
		return nil, errors.New("[ERROR] service name can not be empty")
	}

	host := net.ParseIP(v.GetString("service.host"))
	if host == nil {
		return nil, fmt.Errorf("[ERROR] service host: %s is not a ipv4 address", host)
	}
	return &ServiceConfig{
		ServiceName: v.GetString("service.service_name"),
		Host:        host.String(),
		Port:        v.GetInt("service.port"),
		EnableHTTP:  v.GetBool("service.enable_http"),
	}, nil
}
