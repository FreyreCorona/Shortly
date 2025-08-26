package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/FreyreCorona/Shortly/src/redirect_svc/internal/application"
	"github.com/FreyreCorona/Shortly/src/redirect_svc/internal/infrastructure/cache/valkey"
	httpAdapter "github.com/FreyreCorona/Shortly/src/redirect_svc/internal/infrastructure/http"
)

func main() {
	// Valkey adapter
	cache, err := valkey.NewValkeyCache(fmt.Sprintf("%s:%s", os.Getenv("VALKEY_ADDR"), os.Getenv("VALKEY_PORT")),

		os.Getenv("VALKEY_USERNAME"),
		os.Getenv("VALKEY_PASSWORD"))
	if err != nil {
		log.Fatalf("cache connection error :%s", err.Error())
	}
	// stablish the adapter in the service
	// TODO: REFERENCE THE REPO OBJECT FOR NewRedirectionService PARAMETER
	service := application.NewRedirectionService(cache, nil)
	handler := httpAdapter.NewHandler(service)

	mux := http.NewServeMux()
	handler.Routes(mux)

	// running the service
	runningPort, err := strconv.Atoi(os.Getenv("REDIRECT_SVC_PORT"))
	if err != nil {
		log.Fatal("Uknown port")
	}
	fmt.Printf("Service running on : %d \n", runningPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", runningPort), mux))
}
