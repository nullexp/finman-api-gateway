package main

import (
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	authv1 "github.com/nullexp/finman-gateway-service/grpc/auth/v1"
	"github.com/nullexp/finman-gateway-service/http"
	ginapi "github.com/nullexp/finman-gateway-service/pkg/infrastructure/http/gin"
	"github.com/nullexp/finman-gateway-service/pkg/infrastructure/http/protocol/model"
	"github.com/nullexp/finman-gateway-service/pkg/infrastructure/http/protocol/model/openapi"
	logger "github.com/nullexp/finman-gateway-service/pkg/infrastructure/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("Starting the server")
	logger.Initialize()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	api := ginapi.NewGinApp()

	authUrl := os.Getenv("FINMAN_AUTH_URL")
	port := os.Getenv("PORT")
	ip := os.Getenv("IP")

	api.SetContact(openapi.Contact{Name: "Hope Golestany", Email: "hopegolestany@gmail.com", URL: "https://github.com/nullexp"})
	api.SetInfo(openapi.Info{Version: "0.1", Description: "Api definition for finman", Title: "Finman Api Definition"})
	api.SetLogPolicy(model.LogPolicy{LogBody: false, LogEnabled: false})
	api.SetCors([]string{"http://localhost:8085"})
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := establishGRPCConnection(authUrl, 10)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	authClient := authv1.NewAuthServiceClient(conn)
	auth := http.NewSession(authClient)
	api.AppendModule(auth)

	portValue, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalln(err)
	}

	err = api.EnableOpenApi("/openapi")
	if err != nil {
		log.Fatalln(err)
	}
	err = api.Run(ip, uint(portValue), "debug")
	if err != nil {
		log.Fatalln(err)
	}
}

// establishGRPCConnection establishes a gRPC connection with retry mechanism
func establishGRPCConnection(serverAddr string, retryAttempts int) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error

	for i := 0; i < retryAttempts; i++ {
		conn, err = grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials())) // insecure for test purpose
		if err == nil {
			log.Println("connected")
			return conn, nil
		}
		log.Printf("Failed to connect (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second) // Retry after 2 seconds
	}
	return nil, err
}