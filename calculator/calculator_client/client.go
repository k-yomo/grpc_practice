package main

import (
	"context"
	"fmt"
	"github.com/k-yomo/grpc_practice/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"strconv"
	"time"
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
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	doErrorUnary(c, -2)
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

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Num: int32(10),
		},
		&calculatorpb.ComputeAverageRequest{
			Num: int32(9),
		},
		&calculatorpb.ComputeAverageRequest{
			Num: int32(8),
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComuteAverage: %v", err)
	}

	for _, req := range requests {
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while reciving respose from ComputeAverage: %v", err)
	}
	fmt.Printf("ComputeAverage Response: %v\n", res)
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient)  {
	fmt.Println("Starting to do a BiDi Streaming RPC...")
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		nums := []int32{4, 8, 2, 12}
		for _, num := range nums {
			fmt.Printf("Num: %d\n", num)
			err = stream.Send(&calculatorpb.FindMaximumRequest{
				Num: num,
			})
			if err != nil {
				log.Fatalf("Error while sending request: %v", err)
				time.Sleep(500 * time.Millisecond)
			}
		}
		stream.CloseSend()
	}()
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving response: %v", err)
			}
			fmt.Printf("Response: %v\n", res)
		}
		close(waitc)
	}()
	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient, num int32)  {
	fmt.Println("Starting to do a Sqarearoot Unary RPC...")
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Num: num})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("Negative number was sent")
			}
		} else {
			log.Fatalf("Big Error calling SqareRoot: %v", err)
		}
	}

	fmt.Printf("Result of sqare root of %v:%v\n", num, res.GetNumRoot())
}

