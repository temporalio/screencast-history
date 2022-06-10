package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/hokaccha/go-prettyjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	workflowservice "go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

func NewWorkflowTaskInterceptor() grpc.UnaryClientInterceptor {
	marshaler := jsonpb.Marshaler{Indent: "\t"}

	return func(ctx context.Context, method string, req, response interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		switch o := req.(type) {
		case *workflowservice.PollWorkflowTaskQueueRequest:
			j, _ := marshaler.MarshalToString(o)
			msg, _ := prettyjson.Format([]byte(j))
			fmt.Printf("-> PollWorkflowTaskQueueRequest:\n%s\n\n", msg)
		case *workflowservice.RespondWorkflowTaskCompletedRequest:
			j, _ := marshaler.MarshalToString(o)
			msg, _ := prettyjson.Format([]byte(j))
			fmt.Printf("-> RespondWorkflowTaskCompletedRequest:\n%s\n\n", msg)
		}

		err := invoker(ctx, method, req, response, cc, opts...)
		if err != nil {
			return err
		}

		switch o := response.(type) {
		case *workflowservice.PollWorkflowTaskQueueResponse:
			j, _ := marshaler.MarshalToString(o)
			msg, _ := prettyjson.Format([]byte(j))
			fmt.Printf("<- PollWorkflowTaskQueueResponse:\n%s\n\n", msg)
		}

		return nil
	}
}

func NewActivityTaskInterceptor() grpc.UnaryClientInterceptor {
	marshaler := jsonpb.Marshaler{Indent: "\t"}

	return func(ctx context.Context, method string, req, response interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		switch o := req.(type) {
		case *workflowservice.PollActivityTaskQueueRequest:
			j, _ := marshaler.MarshalToString(o)
			msg, _ := prettyjson.Format([]byte(j))
			fmt.Printf("-> PollActivityTaskQueueRequest:\n%s\n\n", msg)
		case *workflowservice.RespondActivityTaskCompletedRequest:
			j, _ := marshaler.MarshalToString(o)
			msg, _ := prettyjson.Format([]byte(j))
			fmt.Printf("-> RespondActivityTaskCompletedRequest:\n%s\n\n", msg)
		}

		err := invoker(ctx, method, req, response, cc, opts...)
		if err != nil {
			return err
		}

		switch o := response.(type) {
		case *workflowservice.PollActivityTaskQueueResponse:
			j, _ := marshaler.MarshalToString(o)
			msg, _ := prettyjson.Format([]byte(j))
			fmt.Printf("<- PollActivityTaskQueueResponse:\n%s\n\n", msg)
		}

		return nil
	}
}

func main() {
	workflowTaskInterceptor := NewWorkflowTaskInterceptor()

	grpcClient, err := grpc.Dial(
		"127.0.0.1:7233",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(workflowTaskInterceptor),
	)
	defer func() { _ = grpcClient.Close() }()

	if err != nil {
		log.Fatalf("unable to create client: %v", err)
	}

	workflowClient := workflowservice.NewWorkflowServiceClient(grpcClient)

	listener, err := net.Listen("tcp", "127.0.0.1:7234")
	if err != nil {
		log.Fatalf("unable to create listener: %v", err)
	}

	server := grpc.NewServer()
	handler, err := client.NewWorkflowServiceProxyServer(
		client.WorkflowServiceProxyOptions{Client: workflowClient},
	)
	if err != nil {
		log.Fatalf("unable to create service proxy: %v", err)
	}

	workflowservice.RegisterWorkflowServiceServer(server, handler)

	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("unable to serve: %v", err)
	}
}
