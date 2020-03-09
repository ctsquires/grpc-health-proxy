package hello

import (
	"context"
	"log"
)

type Server struct {
}

func NewHelloServer() GreeterServer {
	return &Server{}
}

func (s *Server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &HelloReply{Message: "Hello " + in.GetName()}, nil
}
