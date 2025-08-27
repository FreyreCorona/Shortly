package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/FreyreCorona/Shortly/src/shortener_svc/internal/application"
	db "github.com/FreyreCorona/Shortly/src/shortener_svc/internal/infrastructure/db/postgres"
	"github.com/FreyreCorona/Shortly/src/shortener_svc/internal/infrastructure/grpc"
	httpAdapter "github.com/FreyreCorona/Shortly/src/shortener_svc/internal/infrastructure/http"
)

func main() {
	// postgres adapter
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSGTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))

	repo, err := db.NewPostgresDBRepository(dsn)
	if err != nil {
		log.Fatalf("database connection error :%s", err.Error())
	}
	// stablish the adapter in the service CreateURL
	CreateURLService := application.NewCreateURLService(repo)
	handler := httpAdapter.NewHandler(CreateURLService)

	mux := http.NewServeMux()
	handler.Routes(mux)

	// stablish the adapter in the service RetrieveURL
	// TODO: IMPLEMENT REPO OBJECT
	RetrieveURLService := application.NewRetrieveURLService(repo)
	gRPCServer := grpc.NewGRPCServer(*RetrieveURLService)

	// running the service
	runningPort, err := strconv.Atoi(os.Getenv("URL_SHORTENER_SVC_PORT"))
	if err != nil {
		log.Fatal("Uknown port")
	}
	fmt.Printf("Service running on : %d \n", runningPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", runningPort), mux))
}
