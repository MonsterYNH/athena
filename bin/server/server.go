package main

import (
	"athena/api/v1/health"
	"athena/api/v1/helloworld"
	"athena/registry"
	"athena/server"
	"athena/server/interceptor"
	"athena/services"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthInterceptor("jwt_key")),
	)
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(incomeingHeaderMatcher),
		runtime.WithForwardResponseOption(outgoingHeaderFilter),
	)
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	helloworld.RegisterHellowroldServiceServer(grpcServer, new(services.HelloworldService))
	helloworld.RegisterHellowroldServiceHandlerFromEndpoint(ctx, mux, "localhost:8080", opts)

	health.RegisterHealthServiceServer(grpcServer, new(services.HealthService))
	health.RegisterHealthServiceHandlerFromEndpoint(ctx, mux, "localhost:8080", opts)

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

func incomeingHeaderMatcher(key string) (string, bool) {
	switch key {
	case "ATHENA_TOKEN":
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
		log.Println(values)
	}

	return nil
}
