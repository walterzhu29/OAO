package main

import (
	"context"
	"log"
	"time"

	pb "submitter/protos/model_server" // 更新为实际生成的protobuf包路径

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewModelServerClient(conn)

	request := &pb.ModelRequest{
		RequestId: uuid.New().String(),
		ModelId:   1,
		Input:     []byte("sample input"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.ModelCall(ctx, request)
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}
	log.Printf("Response: %v", response)
}
