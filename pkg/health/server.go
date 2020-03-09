package health

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ctsquires/grpc-health-proxy/pkg/health/healthproto"
)

type Server struct {
	statusMap map[string]pb.HealthCheckResponse_ServingStatus
}

// NewHealthServer creates task service
func NewHealthServer() *Server {
	return &Server{
		statusMap: map[string]pb.HealthCheckResponse_ServingStatus{"": pb.HealthCheckResponse_SERVING},
	}
}

func (s *Server) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	log.Printf("Health Check Hit")
	if servingStatus, ok := s.statusMap[req.Service]; ok {
		return &pb.HealthCheckResponse{
			Status: servingStatus,
		}, nil
	}
	return nil, status.Error(codes.NotFound, "unknown service")
}

func (s *Server) SetServingStatus(service string, servingStatus pb.HealthCheckResponse_ServingStatus) {
	s.statusMap[service] = servingStatus
}
