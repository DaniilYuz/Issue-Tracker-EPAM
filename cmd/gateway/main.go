package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/pkg/gen"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	_ = godotenv.Load()

	grpcPort := os.Getenv("PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	gatewayPort := os.Getenv("GATEWAY_PORT")
	if gatewayPort == "" {
		gatewayPort = "50052"
	}

	grpcHost := os.Getenv("GRPC_HOST")
	if grpcHost == "" {
		grpcHost = "localhost"
	}

	grpcEndpoint := fmt.Sprintf("%s:%s", grpcHost, grpcPort)

	mux := runtime.NewServeMux()

	//set connections to grpc
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := gen.RegisterUserServiceHandlerFromEndpoint(
		context.Background(), mux, grpcEndpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register user service: %v", err)
	}

	err = gen.RegisterIssueServiceHandlerFromEndpoint(
		context.Background(), mux, grpcEndpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register issue service: %v", err)
	}

	err = gen.RegisterProjectServiceHandlerFromEndpoint(
		context.Background(), mux, grpcEndpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register project service: %v", err)
	}

	log.Printf("HTTP Gateway server is running on port %s...", gatewayPort)
	log.Printf("gRPC server should be running on port %s", grpcPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", gatewayPort), mux); err != nil {
		log.Fatalf("Failed to serve HTTP gateway: %v", err)
	}
}
