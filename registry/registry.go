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
	Regist() (interface{}, error)
	Stop()
}
