package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/markustenghamn/beubo/grpc"
	"google.golang.org/grpc"
)

const (
	host = "localhost"
	port = "55051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
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
