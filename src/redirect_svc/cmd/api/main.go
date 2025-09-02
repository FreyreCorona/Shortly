package main

import (
	"context"
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
	"github.com/FreyreCorona/Shortly/src/redirect_svc/internal/infrastructure/rabbitmq"
)

func main() {
	// Valkey adapter
	cache, err := valkey.NewValkeyCache(fmt.Sprintf("%s:%s", os.Getenv("VALKEY_HOST"), os.Getenv("VALKEY_PORT")),
		os.Getenv("VALKEY_USERNAME"),
		os.Getenv("VALKEY_PASSWORD"))
	if err != nil || cache == nil {
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
			log.Fatalf("error on http handler :%v", err)
		}
	})

	address := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_DEFAULT_USER"),
		os.Getenv("RABBITMQ_DEFAULT_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"))

	// stablish listening on message queue
	wg.Go(func() {
		if err := StartQueueConsumer(cache, address); err != nil {
			log.Fatalf("error on queueConsumer :%v", err)
		}
	})
	wg.Wait()
}

func StartHTTPHandler(cache domain.URLCacheRepository, repo domain.URLRepository) error {
	service := application.NewRedirectionService(cache, repo)
	handler := httpAdapter.NewHandler(service)

	mux := http.NewServeMux()
	handler.Routes(mux)

	port := ":" + os.Getenv("REDIRECT_SVC_PORT")

	log.Printf("HTTP handler running on %s", port)

	return http.ListenAndServe(port, mux)
}

func StartQueueConsumer(cache domain.URLCacheRepository, address string) error {
	service := application.NewSetURLService(cache)
	consumer, err := rabbitmq.NewConsumer(*service, address)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Start listening for messages from the queue")
	err = consumer.Listen(ctx)
	if err != nil {
		return err
	}
	defer ctx.Done()
	return nil
}
