package registry

import (
	"google.golang.org/grpc/resolver"
)

type RegistryAble interface {
	// resolver builder
	Scheme() string
	Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error)
	// resolver
	ResolveNow(o resolver.ResolveNowOptions)
	Close()

	// regist
	Regist() error
	Stop() error

	SetConfig(config *RegistryConfig) error
}

type Registry struct {
	Registry RegistryAble
	config   *RegistryConfig
}

func NewRegistry(registry RegistryAble, opts ...RegistryOption) (*Registry, error) {
	config := &RegistryConfig{}

	for _, option := range opts {
		if err := option(config); err != nil {
			return nil, err
		}
	}

	if err := registry.SetConfig(config); err != nil {
		return nil, err
	}

	resolver.Register(registry)

	return &Registry{
		Registry: registry,
		config:   config,
	}, nil
}

func (registry *Registry) Run() error {
	return registry.Registry.Regist()
}

func (registry *Registry) Stop() error {
	return registry.Registry.Stop()
}
