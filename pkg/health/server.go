package health

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ctsquires/grpc-health-proxy/pkg/health/healthpb"
)

type Server struct {
	mu     sync.RWMutex
	ready  bool
	status healthpb.HealthCheckResponse_ServingStatus
}

// NewHealthServer creates task service
func NewHealthServer() *Server {
	return &Server{
		status: healthpb.HealthCheckResponse_SERVING,
	}
}

// Check implements `service Health`.
func (s *Server) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &healthpb.HealthCheckResponse{
		Status: s.status,
	}, nil
}

// Ready implements `service Health`.
func (s *Server) Ready(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.ready {
		return &healthpb.HealthCheckResponse{
			Status: s.status,
		}, nil
	}
	return nil, status.Error(codes.Unavailable, "service not ready")
}

// SetServingStatus is called when need to reset the serving status of a service
func (s *Server) SetServingStatus(servingStatus healthpb.HealthCheckResponse_ServingStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = servingStatus
}

// SetReadyStatus is called to set the service as ready or not ready
func (s *Server) SetReadyStatus(readyStatus bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ready = readyStatus
}
