package health

import (
	"context"
	"fmt"
	"net/http"

	healthpb "github.com/ctsquires/grpc-health-proxy/pkg/grpc_health_proxy"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func HTTPHealthServerFromPort(ctx context.Context, grpcServer *grpc.Server, port int, services []string) (*http.Server, error) {
	healthServer := ConfigureGRPCHealthServer(grpcServer, services)

	mux, err := ConfigureHTTPHealthServer(ctx, healthServer)
	if err != nil {
		return nil, fmt.Errorf("Could not register health handler server ")
	}

	healthHTTPServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	return healthHTTPServer, nil
}

func ConfigureGRPCHealthServer(grpcServer *grpc.Server, services []string) healthpb.HealthServer {
	healthServer := NewHealthServer()
	for _, value := range services {
		healthServer.SetServingStatus(value, healthpb.HealthCheckResponse_SERVING)
	}
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	return healthServer
}

func ConfigureHTTPHealthServer(ctx context.Context, server healthpb.HealthServer) (*http.ServeMux, error) {
	runMux := runtime.NewServeMux()
	if err := healthpb.RegisterHealthHandlerServer(ctx, runMux, server); err != nil {
		return nil, fmt.Errorf("Could not register health handler server ")
	}

	mux := http.NewServeMux()
	mux.Handle("/", runMux)
	return mux, nil
}
