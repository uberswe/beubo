package main

import (
	"context"
	"log"
	"time"

	pb "github.com/markustenghamn/beubo/grpc"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051" // Todo make this a configurable flag
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBeuboGRPCClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Register(ctx, &pb.RegisterRequest{
		Name:       "Plugin",
		Identifier: "Plugin",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
