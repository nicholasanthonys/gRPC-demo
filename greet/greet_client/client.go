package main

import (
	"context"
	"fmt"
	"github.com/nicholasanthonys/gRPC-demo/greet/greetpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	fmt.Println("Hello I'm a client")

	// by default, grpc has ssl. for now, we don't have ssl certificate. remove this in production
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("created client : %v", c)

	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Nicholas",
			LastName:  "Anthony",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC : %v", err)
	}

	log.Printf("response from greet :  %v", res.Result)
}
