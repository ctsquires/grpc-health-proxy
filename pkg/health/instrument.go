package health

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func HTTPHealthServerFromPort(ctx context.Context, grpcServer *grpc.Server, port int, services []string) (*http.Server, error) {
	healthServer := NewHealthServer()
	for _, value := range services {
		healthServer.SetServingStatus(value, HealthCheckResponse_SERVING)
	}
	runMux := runtime.NewServeMux()
	RegisterHealthServer(grpcServer, healthServer)

	if err := RegisterHealthHandlerServer(ctx, runMux, healthServer); err != nil {
		return nil, fmt.Errorf("Could not register health handler server ")
	}

	mux := http.NewServeMux()
	mux.Handle("/", runMux)

	healthHTTPServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	return healthHTTPServer, nil
}
