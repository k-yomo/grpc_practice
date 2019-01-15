package main

import (
	"context"
	"fmt"
	"github.com/yomoda07/grpc_go_course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Recieved Sum RPC: %v", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	sum := firstNumber + secondNumber
	res := &calculatorpb.SumResponse{
		SumResult: sum,
	}
	return res, nil
}

func (*server) Multiply(ctx context.Context, req *calculatorpb.MultiplyRequest) (*calculatorpb.MultiplyResponse, error) {
	fmt.Printf("Recieved Sum RPC: %v", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	res := &calculatorpb.MultiplyResponse{
		MultiplyResult: firstNumber * secondNumber,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Recieved PrimaryNumberDecomposition RPC: %v", req)
	num := req.GetNum()
	divisor := int64(2)
	for num > 1 {
		if num % divisor == 0 {
			num /= divisor
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			}
			err := stream.Send(res)
			if err != nil {
				log.Fatalf("Error while sending response: %v", err)
			}
		} else {
			divisor++
		}
	}
	return nil
}

func main() {
	fmt.Println("Start server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}