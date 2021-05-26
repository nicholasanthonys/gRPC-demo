package main

import (
	"context"
	"fmt"
	"github.com/nicholasanthonys/gRPC-demo/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
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

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDiStreaming(c)
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

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do server streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Nicholas",
			LastName:  "Anthony",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC : %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// reach the end of the stream
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}

		log.Printf("repsonse from greet many times : \n %v ", msg.GetResult())

	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do sclient streaming RPC...")
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling long greet %v", err)
		return
	}
	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Nicholas",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Anthony",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Suhartono",
			},
		},
	}

	for _, req := range requests {
		fmt.Printf("sending request : %v \n", req)
		err := stream.Send(req)
		time.Sleep(100 * time.Millisecond)
		if err != nil {
			return
		}
	}

	// return longGreet response and error
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from long greet %v", err)
		return
	}

	fmt.Printf("Long greet response %v\n : ", res)

}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do BiDi streaming RPC...")

	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Nicholas",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Anthony",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Suhartono",
			},
		},
	}

	waitc := make(chan struct{})

	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message %v\n", req)
			err := stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
			if err != nil {
				fmt.Printf("Error while send streaming %v", err)
				return
			}
		}
		err := stream.CloseSend()
		if err != nil {
			fmt.Printf("Error while close send streaming %v", err)
			return
		}

	}()

	// we receive a bunch of messages to the client (go routine)
	go func() {

		for {
			// function to receive a bunch of messages
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				fmt.Printf("error while receiving %v", err)
				break
			}

			fmt.Printf("Receive: %v\n", res)
		}
		close(waitc)

	}()

	// block until everything is done
	<-waitc
}
