package main

import (
	"context"
	"fmt"
	"github.com/yomoda07/grpc_go_course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"strconv"
)

func main() {
	fmt.Println("Calculator Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//doUnary(c)
	doServerStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &calculatorpb.MultiplyRequest{
		FirstNumber: 5,
		SecondNumber: 10,
	}

	res, err := c.Multiply(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.MultiplyResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient)  {
	fmt.Println("Starting to do a Server Streaming RPC...")
	num := int64(12390392)
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Num: num,
	}
	var result []string
	result = append(result, strconv.FormatInt(num, 10), " = ")

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecompostion RPC: %V", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		if len(result) > 2 {
			result = append(result, " * ")
		}
		result = append(result, strconv.FormatInt(msg.GetPrimeFactor(), 10))
	}
	log.Println(result)
}