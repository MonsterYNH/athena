package registry

import "time"

type RegistryConfig struct {
	Name           string
	IP             string
	TTL            int64
	Entrypoints    []string
	RouteInfos     []string
	DialTimeout    time.Duration
	DependServices []string
}

type RegistryOption func(*RegistryConfig) error

func WithEntrypointsRegistryOption(entrypoints []string) RegistryOption {
	return func(config *RegistryConfig) error {
		config.Entrypoints = entrypoints
		return nil
	}
}

func WithNameRegistryOption(name string) RegistryOption {
	return func(config *RegistryConfig) error {
		config.Name = name
		return nil
	}
}

func WithIPRegistryOption(ip string) RegistryOption {
	return func(config *RegistryConfig) error {
		config.IP = ip
		return nil
	}
}

func WithDialTimeoutRegistryOption(dur time.Duration) RegistryOption {
	return func(config *RegistryConfig) error {
		config.DialTimeout = dur
		return nil
	}
}

func WithTTLRegistryOption(ttl int64) RegistryOption {
	return func(config *RegistryConfig) error {
		config.TTL = ttl
		return nil
	}
}

func WithRouteInfosOption(routeInfos []string) RegistryOption {
	return func(config *RegistryConfig) error {
		config.RouteInfos = routeInfos
		return nil
	}
}

func WithDependServices(services []string) RegistryOption {
	return func(config *RegistryConfig) error {
		config.DependServices = services
		return nil
	}
}
