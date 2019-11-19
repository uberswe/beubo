package beubo

import (
	"context"
	pb "github.com/markustenghamn/beubo/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

// Grpc

//go:generate protoc -I grpc --go_out=plugins=grpc:grpc grpc/beubo.proto

const (
	grpcPort = ":50051"
)

type server struct{}

func (s *server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.RegisterResponse{Message: "Registered", Success: true}, nil
}

func (s *server) Insert(ctx context.Context, in *pb.InsertRequest) (*pb.InsertResponse, error) {
	log.Printf("Received: ")
	return &pb.InsertResponse{}, nil
}

func (s *server) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	log.Printf("Received: ")
	return &pb.UpdateResponse{}, nil
}

func (s *server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	log.Printf("Received: ")
	return &pb.DeleteResponse{}, nil
}

func (s *server) Select(ctx context.Context, in *pb.SelectRequest) (*pb.SelectResponse, error) {
	log.Printf("Received: ")
	return &pb.SelectResponse{}, nil
}

func (s *server) Handle(ctx context.Context, in *pb.HandleRequest) (*pb.HandleResponse, error) {
	log.Printf("Received: ")
	return &pb.HandleResponse{}, nil
}

func (s *server) FetchEndpoints(ctx context.Context, in *pb.EmptyRequest) (*pb.FetchEndpointsResponse, error) {
	log.Printf("Received: ")
	return &pb.FetchEndpointsResponse{}, nil
}

func (s *server) CallEndpoint(ctx context.Context, in *pb.CallEndpointRequest) (*pb.CallEndpointResponse, error) {
	log.Printf("Received: ")
	return &pb.CallEndpointResponse{}, nil
}

func grpcInit() {
	log.Printf("Starting grpc server")
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %s", grpcPort)
	s := grpc.NewServer()
	pb.RegisterBeuboGRPCServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
