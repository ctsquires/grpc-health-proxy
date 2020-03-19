package health

import (
	"context"
	"log"
	"sync"

	healthpb "github.com/ctsquires/grpc-health-proxy/pkg/health/grpc_health_proxy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	mu        sync.RWMutex
	shutdown  bool
	statusMap map[string]healthpb.HealthCheckResponse_ServingStatus
}

// NewHealthServer creates task service
func NewHealthServer() *Server {
	return &Server{
		statusMap: map[string]healthpb.HealthCheckResponse_ServingStatus{"": healthpb.HealthCheckResponse_SERVING},
	}
}

func (s *Server) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	log.Printf("Health Check Hit")
	if servingStatus, ok := s.statusMap[req.Service]; ok {
		return &healthpb.HealthCheckResponse{
			Status: servingStatus,
		}, nil
	}
	return nil, status.Error(codes.NotFound, "unknown service")
}

// SetServingStatus is called when need to reset the serving status of a service
// or insert a new service entry into the statusMap.
func (s *Server) SetServingStatus(service string, servingStatus healthpb.HealthCheckResponse_ServingStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.shutdown {
		log.Printf("health: status changing for %s to %v is ignored because health service is shutdown", service, servingStatus)
		return
	}
	s.statusMap[service] = servingStatus
}

// Shutdown sets all serving status to NOT_SERVING, and configures the server to
// ignore all future status changes.
//
// This changes serving status for all services. To set status for a particular
// services, call SetServingStatus().
func (s *Server) Shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.shutdown = true
	for service := range s.statusMap {
		s.SetServingStatus(service, healthpb.HealthCheckResponse_NOT_SERVING)
	}
}

// Resume sets all serving status to SERVING, and configures the server to
// accept all future status changes.
//
// This changes serving status for all services. To set status for a particular
// services, call SetServingStatus().
func (s *Server) Resume() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.shutdown = false
	for service := range s.statusMap {
		s.SetServingStatus(service, healthpb.HealthCheckResponse_SERVING)
	}
}
