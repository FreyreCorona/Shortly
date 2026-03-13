package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	// stablish the adapter in the service
	wg.Go(func() {
		if err := StartHTTPHandler(ctx, cache, repo); err != nil {
			log.Printf("error on http handler :%v", err)
		}
	})

	address := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_DEFAULT_USER"),
		os.Getenv("RABBITMQ_DEFAULT_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"))

	// stablish listening on message queue
	wg.Go(func() {
		if err := StartQueueConsumer(ctx, cache, address); err != nil {
			log.Printf("error on queueConsumer :%v", err)
		}
	})

	<-ctx.Done()
	log.Println("Shutting down gracefully...")
	wg.Wait()
	log.Println("Server stopped")
}

func StartHTTPHandler(ctx context.Context, cache domain.URLCacheRepository, repo domain.URLRepository) error {
	service := application.NewRedirectionService(cache, repo)
	handler := httpAdapter.NewHandler(service)

	mux := http.NewServeMux()
	handler.Routes(mux)

	port := ":" + os.Getenv("REDIRECT_SVC_PORT")
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("HTTP handler running on %s", port)

	go func() {
		<-ctx.Done()
		log.Println("Stopping HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP shutdown error: %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func StartQueueConsumer(ctx context.Context, cache domain.URLCacheRepository, address string) error {
	service := application.NewSetURLService(cache)
	consumer, err := rabbitmq.NewConsumer(*service, address)
	if err != nil {
		return err
	}

	log.Println("Start listening for messages from the queue")
	err = consumer.Listen(ctx)
	if err != nil {
		return err
	}
	return nil
}
