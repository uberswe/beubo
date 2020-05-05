package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	pb "github.com/markustenghamn/beubo/grpc"
	"google.golang.org/grpc"
	"io"
	"log"
)

const (
	pluginName       = "Custom Page Example Plugin"
	pluginIdentifier = "custom-page-example"
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
	stream, err := c.Connect(ctx)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	for {
		m := pb.PluginMessage{
			Name:       pluginName,
			Identifier: pluginIdentifier,
		}

		serialized, err := proto.Marshal(&m)
		if err != nil {
			log.Println("Could not serialize plugin message")
			return
		}

		err = stream.Send(&pb.Event{
			Key:  "register",
			Data: "custom-page-plugin",
			Values: []*any.Any{
				{
					TypeUrl: proto.MessageName(&m),
					Value:   serialized,
				},
			},
		})

		if err != nil {
			log.Fatalf("could not send: %v", err)
		}

		for {
			request, err := stream.Recv()

			if err == io.EOF {
				// read done.
				return
			}

			if err != nil {
				log.Fatalf("could not receive: %v", err)
			}

			log.Printf("Request received: %s\n", request.Key)

			for _, anyVar := range request.Values {
				log.Println(anyVar.TypeUrl)
				if anyVar.TypeUrl == "beubo.Request" {
					var m pb.Request
					err := proto.Unmarshal(anyVar.Value, &m)
					if err != nil {
						log.Println(err)
						return
					}
					log.Printf("Request message unmarshalled: %s\n", m.Url)
				} else if anyVar.TypeUrl == "beubo.Response" {
					var m pb.Response
					err := proto.Unmarshal(anyVar.Value, &m)
					if err != nil {
						log.Println(err)
						return
					}
					log.Printf("Response message unmarshalled: %s\n", m.Content)
				}
			}
		}
	}
}
