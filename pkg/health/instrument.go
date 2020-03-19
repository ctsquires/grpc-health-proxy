package health

import (
	"context"
	"fmt"
	"net/http"

	healthpb "github.com/ctsquires/grpc-health-proxy/pkg/health/grpc_health_proxy"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func HTTPHealthServerFromPort(ctx context.Context, grpcServer *grpc.Server, port int, services []string) (*http.Server, error) {
	healthServer := NewHealthServer()
	for _, value := range services {
		healthServer.SetServingStatus(value, healthpb.HealthCheckResponse_SERVING)
	}
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	runMux := runtime.NewServeMux()
	if err := healthpb.RegisterHealthHandlerServer(ctx, runMux, healthServer); err != nil {
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
