package beubo

import (
	pb "github.com/markustenghamn/beubo/grpc"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

// Grpc

// protoc -I grpc --go_out=plugins=grpc:grpc grpc/beubo.proto

const (
	grpcPort = ":50051"
)

type server struct{}

func (s *server) Connect(stream pb.BeuboGRPC_ConnectServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Received: %v", in.Name)
	}
}

func (s *server) Requests(pluginMessage *pb.PluginMessage, stream pb.BeuboGRPC_RequestsServer) error {
	for {
		request := <-requestChannel
		if err := stream.Send(&request); err != nil {
			return err
		}
	}
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
