package test

import (
	"athena/api/v1/helloworld"
	"athena/registry"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestRegistry(t *testing.T) {
	_, err := registry.NewRegistry(
		&registry.EtcdRegistry{},
		registry.WithDependServices([]string{"my_test"}),
		registry.WithDialTimeoutRegistryOption(time.Second*5),
		registry.WithEntrypointsRegistryOption([]string{"localhost:2379"}),
		registry.WithIPRegistryOption("localhost:8080"),
		registry.WithTTLRegistryOption(5),
	)
	if err != nil {
		panic(err)
	}

	// resolver.Register(reg.Registry)
	conn, err := grpc.Dial("etcd://health.HealthService/HealthCheck", grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))
	if err != nil {
		t.Fatal(err)
	}

	client := helloworld.NewHellowroldServiceClient(conn)

	resp, err := client.SyaHello(context.Background(), &helloworld.HelloworldRequest{
		Name: "zhangsan",
	})

	if err != nil {
		t.Fatal(err)
	}

	bytes, _ := json.Marshal(resp)
	fmt.Println(string(bytes))
}
