package ops

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	statusMap map[string]OpsListResponse_ServingStatus
	version   string
	config    Config
}

// NewOpsServer creates task service
func NewOpsServer() *Server {
	return &Server{
		statusMap: map[string]OpsListResponse_ServingStatus{"": OpsListResponse_SERVING},
		version:   "dev",
	}
}

func (s *Server) Ops(ctx context.Context, req *OpsListRequest) (*OpsListResponse, error) {
	log.Printf("Operations Ops Hit")
	if servingStatus, ok := s.statusMap[req.Service]; ok {
		return &OpsListResponse{
			Status:  servingStatus,
			Version: s.version,
			Config:  s.config.serveJSON(),
		}, nil
	}
	return nil, status.Error(codes.NotFound, "unknown service")
}

func (s *Server) SetServingStatus(service string, servingStatus OpsListResponse_ServingStatus) {
	s.statusMap[service] = servingStatus
}
