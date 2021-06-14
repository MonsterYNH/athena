package main

import (
	"athena/api/v1/auth2"
	"athena/api/v1/health"
	"athena/registry"
	"athena/server"
	"athena/services"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(headerMatcher),
		runtime.WithForwardResponseOption(outgoingHeaderFilter),
	)
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	health.RegisterHealthServiceServer(grpcServer, new(services.HealthService))
	health.RegisterHealthServiceHandlerFromEndpoint(ctx, mux, "localhost:8080", opts)

	auth2.RegisterAuth2SerivceServer(grpcServer, new(services.Auth2Service))
	auth2.RegisterAuth2SerivceHandlerFromEndpoint(ctx, mux, "localhost:8080", opts)

	engine, err := server.NewServer(grpcServer, mux)
	if err != nil {
		panic(err)
	}

	reg, err := registry.NewRegistry(
		&registry.EtcdRegistry{},
		registry.WithDependServices([]string{"my_test"}),
		registry.WithDialTimeoutRegistryOption(time.Second*5),
		registry.WithEntrypointsRegistryOption([]string{"localhost:2379"}),
		registry.WithIPRegistryOption("localhost:8080"),
		registry.WithNameRegistryOption("my_test"),
		registry.WithTTLRegistryOption(5),
		registry.WithRouteInfosOption(engine.GetGrpcRouteInfo()),
	)
	if err != nil {
		panic(err)
	}

	if err := reg.Run(); err != nil {
		panic(err)
	}
	defer reg.Stop()

	if err := engine.Run(func(sc *server.ServerConfig) error {
		sc.EnableHTTP = true
		sc.Host = "0.0.0.0"
		sc.Port = 8080
		return nil
	}); err != nil {
		panic(err)
	}
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
