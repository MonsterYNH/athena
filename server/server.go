package server

import (
	"athena/registry"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type ServerConfig struct {
	Host       string
	Port       int
	EnableHTTP bool
}

type Server struct {
	registry   *registry.RegistryAble
	grpcServer *grpc.Server
	gwMux      *runtime.ServeMux
}

type ServerConfigOption func(*ServerConfig) error

func NewServer(server *grpc.Server, gw *runtime.ServeMux) (*Server, error) {
	return &Server{
		grpcServer: server,
		gwMux:      gw,
	}, nil
}

func (server *Server) Run(options ...ServerConfigOption) error {
	config := new(ServerConfig)

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

func GetGrpcRouteInfo(server *grpc.Server) []string {
	routeInfos := make([]string, 0)

	for name, info := range server.GetServiceInfo() {
		for _, method := range info.Methods {
			routeInfos = append(routeInfos, fmt.Sprintf("%s/%s", name, method.Name))
		}
	}

	sort.Strings(routeInfos)

	return routeInfos
}
