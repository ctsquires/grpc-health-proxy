package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/ctsquires/grpc-health-proxy/pkg/health"
	"github.com/ctsquires/grpc-health-proxy/pkg/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	appPort    = flag.Int("port", 8081, "The server port")
	healthPort = flag.Int("health-port", 8082, "The server port")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	grpcServer := grpc.NewServer()
	hello.RegisterGreeterServer(grpcServer, hello.NewHelloServer())
	reflection.Register(grpcServer)

	healthHTTPServer, healthGRPCServer, err := health.HTTPHealthServerFromPort(ctx, grpcServer, *healthPort)
	if err != nil {
		log.Fatal("Could not register health handler server")
	}

	healthGRPCServer.SetReadyStatus(true)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *appPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	errChan := make(chan error)
	go func() {
		log.Println("HelloWorld gRPC Server Listening on", *appPort)
		errChan <- grpcServer.Serve(lis)
	}()
	go func() {
		log.Println("Health HTTP Server Listening On", *healthPort)
		errChan <- healthHTTPServer.ListenAndServe()
	}()

	log.Fatal(ctx, <-errChan)
}
