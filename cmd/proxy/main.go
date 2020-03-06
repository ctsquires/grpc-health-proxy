package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/ctsquires/grpc-health-proxy/pkg/health"
	"github.com/ctsquires/grpc-health-proxy/pkg/health/healthproto"
	"github.com/ctsquires/grpc-health-proxy/pkg/hello"
	"github.com/ctsquires/grpc-health-proxy/pkg/hello/helloproto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	appPort    = flag.Int("port", 8080, "The server port")
	healthPort = flag.Int("health-port", 8082, "The server port")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	grpcServer := grpc.NewServer()
	helloproto.RegisterGreeterServer(grpcServer, hello.NewHelloServer())
	reflection.Register(grpcServer)

	healthServer := health.NewHealthServer()
	runMux := runtime.NewServeMux()
	healthproto.RegisterHealthServer(grpcServer, healthServer)
	if err := healthproto.RegisterHealthHandlerServer(ctx, runMux, healthServer); err != nil {
		log.Fatal("Could not register health handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", runMux)
	healthHttpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *healthPort),
		Handler: mux,
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
		log.Println(ctx, "ops listening on port", *healthPort)
		errChan <- healthHttpServer.ListenAndServe()
	}()

	// Start the server
	log.Fatal(ctx, <-errChan)
}
