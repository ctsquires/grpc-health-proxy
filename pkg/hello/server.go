package hello

import (
	"context"
	"log"

	pb "github.com/ctsquires/grpc-health-proxy/pkg/hello/helloproto"
)

type Server struct {
}

func NewHelloServer() pb.GreeterServer {
	return &Server{}
}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
