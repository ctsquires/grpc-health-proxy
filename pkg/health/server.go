package health

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	statusMap map[string]HealthCheckResponse_ServingStatus
}

// NewHealthServer creates task service
func NewHealthServer() *Server {
	return &Server{
		statusMap: map[string]HealthCheckResponse_ServingStatus{"": HealthCheckResponse_SERVING},
	}
}

func (s *Server) Check(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	log.Printf("Health Check Hit")
	if servingStatus, ok := s.statusMap[req.Service]; ok {
		return &HealthCheckResponse{
			Status: servingStatus,
		}, nil
	}
	return nil, status.Error(codes.NotFound, "unknown service")
}

func (s *Server) SetServingStatus(service string, servingStatus HealthCheckResponse_ServingStatus) {
	s.statusMap[service] = servingStatus
}
