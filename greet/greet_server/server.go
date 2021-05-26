package main

import (
	"context"
	"fmt"
	"github.com/nicholasanthonys/gRPC-demo/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked. request is \n %v", req)

	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest,
	stream greetpb.GreetService_GreetManyTimesServer) error {

	fmt.Printf("GreetManyTimes function was invoked with %v\n", req)
	firstName := req.Greeting.GetFirstName()

	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: "Hello " + firstName + " number : " + strconv.Itoa(i),
		}

		err := stream.Send(res)
		if err != nil {
			return err
		}

		// sleep 1 second
		time.Sleep(1 * time.Second)
	}
	return nil
}
func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked with a streaming request %n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName

		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})

		if err != nil {
			log.Fatalf("Error while sending data to client %v", err)
			return err
		}
	}
}
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function was invoked with streaming request")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "!"
	}
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Fai;led to listen : %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : %v ", err)
	}
}
