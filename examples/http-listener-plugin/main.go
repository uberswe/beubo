package main

import (
	"context"
	"fmt"
	pb "github.com/markustenghamn/beubo/grpc"
	"google.golang.org/grpc"
	"io"
	"log"
)

const (
	pluginName       = "HTTP Listener Example Plugin"
	pluginIdentifier = "http-listener-example"
	remoteHost       = "localhost"
	remotePort       = "50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", remoteHost, remotePort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Printf("Connected to %s:%s\n", remoteHost, remotePort)
	defer conn.Close()
	c := pb.NewBeuboGRPCClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, err := c.Requests(ctx, &pb.PluginMessage{
		Name:       pluginName,
		Identifier: pluginIdentifier,
	})
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	for {
		request, err := r.Recv()

		if err == io.EOF {
			// read done.
			return
		}

		if err != nil {
			log.Fatalf("could not receive: %v", err)
		}

		log.Printf("Request received %s: %s\n", request.Url, request.Method)
		for _, header := range request.Headers {
			log.Printf("%s: %s\n", header.Key, header.Values[0])
		}
	}
}
