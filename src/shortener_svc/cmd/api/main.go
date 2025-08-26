package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/FreyreCorona/Shortly/src/shortener_svc/internal/application"
	db "github.com/FreyreCorona/Shortly/src/shortener_svc/internal/infrastructure/db/postgres"
	httpAdapter "github.com/FreyreCorona/Shortly/src/shortener_svc/internal/infrastructure/http"
)

func main() {
	// postgres addapter
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSGTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))

	repo, err := db.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("database connection error :%s", err.Error())
	}
	service := application.NewURLService(repo)
	handler := httpAdapter.NewHandler(service)

	mux := http.NewServeMux()
	handler.Routes(mux)
	runningPort, err := strconv.Atoi(os.Getenv("URL_SHORTENER_SVC_PORT"))
	if err != nil {
		log.Fatal("Uknown port")
	}
	fmt.Printf("Server running on : %d \n", runningPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", runningPort), mux))
}
