package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/MonsterYNH/athena/registry"

	"github.com/MonsterYNH/athena/config"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type Server struct {
	registry        *registry.RegistryAble
	grpcServer      *grpc.Server
	gwMux           *runtime.ServeMux
	authToken       string
	enableAuthCheck bool
}

type ServerConfigOption func(*config.ServiceConfig) error

func New(serverOpts []grpc.ServerOption, serveMuxOpts []runtime.ServeMuxOption) (*grpc.Server, *runtime.ServeMux) {
	grpcServer := grpc.NewServer(serverOpts)

	serveMuxOpts = append(serveMuxOpts, runtime.WithIncomingHeaderMatcher(headerMatcher))
	serveMuxOpts = append(serveMuxOpts, runtime.WithForwardResponseOption(outgoingHeaderFilter))
	mux := runtime.NewServeMux(serveMuxOpts...)

	return grpcServer, mux
}

func NewServer(server *grpc.Server, gw *runtime.ServeMux) (*Server, error) {
	return &Server{
		grpcServer: server,
		gwMux:      gw,
	}, nil
}

func (server *Server) Run(options ...ServerConfigOption) error {
	config := new(config.ServiceConfig)

	for _, option := range options {
		if err := option(config); err != nil {
			return err
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/", server.gwMux)

	log.Println(fmt.Sprintf("grpc service start at %s:%d, http gateway status: %t", config.Host, config.Port, config.EnableHTTP))
	return http.ListenAndServe(
		fmt.Sprintf("%s:%d", config.Host, config.Port),
		h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.EnableHTTP {
				if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
					server.grpcServer.ServeHTTP(w, r)
				} else {
					mux.ServeHTTP(w, r)
				}
			} else {
				server.grpcServer.ServeHTTP(w, r)
			}

		}), &http2.Server{}),
	)
}

func (server *Server) GetGrpcRouteInfo() []string {
	routeInfos := make([]string, 0)

	for name, info := range server.grpcServer.GetServiceInfo() {
		for _, method := range info.Methods {
			routeInfos = append(routeInfos, fmt.Sprintf("%s/%s", name, method.Name))
		}
	}

	sort.Strings(routeInfos)

	return routeInfos
}

func headerMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "athena_token":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func outgoingHeaderFilter(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	if values := md.HeaderMD.Get("athena_token"); len(values) > 0 {
		w.Header().Set("ATHENA_TOKEN", values[0])
	}

	return nil
}
