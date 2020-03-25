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

	//healthGRPCServer.SetReadyStatus(true)
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

	// dummy go routine to set the ready status when a function returns true
	go func() {
		urls := []string{
			"http://localhost:8082",
			"http://localhost:8084",
		}
		var checked bool
		for !checked {
			results, ready := checkDependents(urls)
			if !ready {
				log.Println(results)
			} else {
				checked = true
			}
			healthGRPCServer.SetReadyStatus(ready)
			time.Sleep(1 * time.Second)
		}
	}()

	log.Fatal(ctx, <-errChan)
}

func checkDependents(urls []string) (string, bool) {
	client := &http.Client{}
	var result string
	ok := true
	for _, url := range urls {
		resp, err := client.Get(fmt.Sprintf("%s/healthz", url))
		if err != nil {
			result += fmt.Sprintf("error health checking %s, err: %v \n", url, err)
			ok = false
			continue
		}
		if resp.StatusCode != http.StatusOK {
			result += fmt.Sprintf("error health checking %s, statuscode returned %d \n", url, resp.StatusCode)
			ok = false
		}
	}
	return result, ok
}
