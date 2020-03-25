package health

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ctsquires/grpc-health-proxy/pkg/health/healthpb"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func HTTPHealthServerFromPort(ctx context.Context, grpcServer *grpc.Server, port int) (*http.Server, *Server, error) {
	healthServer := ConfigureGRPCHealthServer(grpcServer)

	mux, err := ConfigureHTTPHealthServer(ctx, healthServer)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not register health handler server ")
	}

	healthHTTPServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	return healthHTTPServer, healthServer, nil
}

func ConfigureGRPCHealthServer(grpcServer *grpc.Server) *Server {
	healthServer := NewHealthServer()
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
