package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"time"

	pb "submitter/protos/model_server"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedModelServerServer
}

func (s *server) ModelCall(ctx context.Context, req *pb.ModelRequest) (*pb.ModelResponse, error) {
	log.Printf("Received request: %v", req)
	time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
	return &pb.ModelResponse{
		RequestId: req.GetRequestId(),
		Message:   "Mock model answer",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterModelServerServer(s, &server{})
	log.Printf("Server is listening on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
