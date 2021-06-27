package config

import (
	"errors"
	"time"

	"github.com/spf13/viper"
)

type RegistryConfig struct {
	Name           string
	IP             string
	TTL            int64
	Entrypoints    []string
	RouteInfos     []string
	DialTimeout    time.Duration
	DependServices []string
}

func newRegistryConfig(v *viper.Viper) (*RegistryConfig, error) {
	name := v.GetString("registry.name")
	if len(name) == 0 {
		return nil, errors.New("[ERROR] registry.name can not be empty")
	}
	ip := v.GetString("registry.ip")
	if len(ip) == 0 {
		return nil, errors.New("[ERROR] registry.ip can not be empty")
	}
	ttl := v.GetInt64("registry.ttl")
	if ttl == 0 {
		return nil, errors.New("[ERROR] registry.ttl can not be empty")
	}
	entrypoints := v.GetStringSlice("registry.entrypoints")
	if len(entrypoints) == 0 {
		return nil, errors.New("[ERROR] registry.entrypoints can not be empty")
	}
	dialTimeout := v.GetInt64("registry.dial_timeout")
	if dialTimeout == 0 {
		return nil, errors.New("[ERROR] registry.dial_timeout can not be empty")
	}
	dependServices := v.GetStringSlice("registry.depend_services")
	if len(dependServices) == 0 {
		return nil, errors.New("[ERROR] registry depend services can not be empty")
	}

	return &RegistryConfig{
		Name:           name,
		IP:             ip,
		TTL:            ttl,
		Entrypoints:    entrypoints,
		DialTimeout:    time.Duration(dialTimeout),
		DependServices: dependServices,
	}, nil
}
