package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/ctsquires/grpc-health-proxy/pkg/health"
	"github.com/ctsquires/grpc-health-proxy/pkg/hello"
	"github.com/ctsquires/grpc-health-proxy/pkg/ops"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	version    = "dev"
	appPort    = flag.Int("port", 8080, "The server port")
	healthPort = flag.Int("health-port", 8082, "The server port")
	opsPort    = flag.Int("ops-port", 8084, "The server port")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	grpcServer := grpc.NewServer()
	hello.RegisterGreeterServer(grpcServer, hello.NewHelloServer())
	reflection.Register(grpcServer)

	opsConfig := &ops.Config{
		Port: *opsPort,
	}
	serviceNames := []string{"helloproto.Greeter"}

	healthHTTPServer, err := health.HTTPHealthServerFromPort(ctx, grpcServer, *healthPort, serviceNames)
	if err != nil {
		log.Fatal("Could not register health handler server")
	}

	opsHTTPServer, err := ops.HTTPOpsServerFromConfig(ctx, grpcServer, opsConfig, serviceNames, version)
	if err != nil {
		log.Fatal("Could not register health handler server")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *appPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	errChan := make(chan error)
	go func() {
		log.Println(ctx, "HelloWorld gRPC Server Listening on", *appPort)
		errChan <- grpcServer.Serve(lis)
	}()
	go func() {
		log.Println(ctx, "Health Http Server Listening On", *healthPort)
		errChan <- healthHTTPServer.ListenAndServe()
	}()
	go func() {
		log.Println(ctx, "Ops Http Server Listening On", 8084)
		errChan <- opsHTTPServer.ListenAndServe()
	}()

	// Start the server
	log.Fatal(ctx, <-errChan)
}
