package main

import (
	"context"
	"fmt"
	"github.com/nicholasanthonys/gRPC-demo/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf(" Request from client is \n %v", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber

	return &calculatorpb.SumResponse{
		SumResult: firstNumber + secondNumber,
	}, nil
}
func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Fai;led to listen : %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : %v ", err)
	}
}
