package health

import (
	"context"
	"log"

	pb "github.com/ctsquires/grpc-health-proxy/pkg/health/healthproto"
)

type Server struct {
}

// NewHealthServer creates task service
func NewHealthServer() pb.HealthServer {
	return &Server{}
}

func (s *Server) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	log.Printf("Health Check Hit")
	return &pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING}, nil
}
