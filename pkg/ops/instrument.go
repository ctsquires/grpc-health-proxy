package ops

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type Config struct {
	Port int
}

func (c *Config) serveJSON() string {
	payload, _ := json.Marshal(c)
	return string(payload)
}

func HTTPOpsServerFromConfig(ctx context.Context, grpcServer *grpc.Server, cfg *Config, services []string, version string) (*http.Server, error) {
	OpsServer := NewOpsServer()

	OpsServer.version = version
	for _, value := range services {
		OpsServer.SetServingStatus(value, OpsListResponse_SERVING)
	}

	runMux := runtime.NewServeMux()
	RegisterOperationsServer(grpcServer, OpsServer)

	if err := RegisterOperationsHandlerServer(ctx, runMux, OpsServer); err != nil {
		return nil, fmt.Errorf("Could not register Ops handler server ")
	}

	mux := http.NewServeMux()
	mux.Handle("/", runMux)
	OpsHTTPServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}
	return OpsHTTPServer, nil
}
