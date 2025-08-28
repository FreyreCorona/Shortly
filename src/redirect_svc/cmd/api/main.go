package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/FreyreCorona/Shortly/src/redirect_svc/internal/application"
	"github.com/FreyreCorona/Shortly/src/redirect_svc/internal/domain"
	"github.com/FreyreCorona/Shortly/src/redirect_svc/internal/infrastructure/cache/valkey"
	"github.com/FreyreCorona/Shortly/src/redirect_svc/internal/infrastructure/grpc"
	httpAdapter "github.com/FreyreCorona/Shortly/src/redirect_svc/internal/infrastructure/http"
)

func main() {
	// Valkey adapter
	cache, err := valkey.NewValkeyCache(fmt.Sprintf("%s:%s", os.Getenv("VALKEY_HOST"), os.Getenv("VALKEY_PORT")),

		os.Getenv("VALKEY_USERNAME"),
		os.Getenv("VALKEY_PASSWORD"))
	if err != nil {
		log.Fatalf("cache connection error :%s", err.Error())
	}

	repo, err := grpc.NewGRPCRepository(fmt.Sprintf("%s:%s", os.Getenv("URL_SHORTENER_SVC_HOST"), os.Getenv("URL_SHORTENER_SVC_GRPC_PORT")))
	if err != nil {
		log.Fatalf("grpc client error :%v", err.Error())
	}

	var wg sync.WaitGroup

	// stablish the adapter in the service
	wg.Go(func() {
		if err := StartHTTPHandler(cache, repo); err != nil {
			log.Printf("error on http handler :%v", err)
		}
	})

	wg.Wait()
}

func StartHTTPHandler(cache domain.URLCacheRepository, repo domain.URLRepository) error {
	service := application.NewRedirectionService(cache, nil)
	handler := httpAdapter.NewHandler(service)

	mux := http.NewServeMux()
	handler.Routes(mux)

	port := ":" + os.Getenv("REDIRECT_SVC_PORT")

	log.Printf("HTTP handler running on %s", port)

	return http.ListenAndServe(port, mux)
}
