package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/FreyreCorona/Shortly/protos"
	"github.com/FreyreCorona/Shortly/src/shortener_svc/internal/application"
	"github.com/FreyreCorona/Shortly/src/shortener_svc/internal/domain"
	db "github.com/FreyreCorona/Shortly/src/shortener_svc/internal/infrastructure/db/postgres"
	grpcadapter "github.com/FreyreCorona/Shortly/src/shortener_svc/internal/infrastructure/grpc"
	httpAdapter "github.com/FreyreCorona/Shortly/src/shortener_svc/internal/infrastructure/http"
	"github.com/FreyreCorona/Shortly/src/shortener_svc/internal/infrastructure/rabbitmq"
	"google.golang.org/grpc"
)

func main() {
	// postgres adapter
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))

	repo, err := db.NewPostgresDBRepository(dsn)
	if err != nil {
		log.Fatalf("database connection error :%s", err.Error())
	}

	publisher, err := rabbitmq.NewProducerPublisher(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_DEFAULT_USER"),
		os.Getenv("RABBITMQ_DEFAULT_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT")))
	if err != nil {
		log.Fatalf("error on rabbitmq producer :%v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	// stablish the adapter in the service CreateURL
	wg.Go(func() {
		if err := StartGRPCServer(ctx, repo); err != nil {
			log.Printf("error on GRPC server :%v", err)
		}
	})

	// stablish the adapter in the service RetrieveURL
	wg.Go(func() {
		if err := StartHTTPHandler(ctx, repo, publisher); err != nil {
			log.Printf("error on HTTP handler :%v", err)
		}
	})

	<-ctx.Done()
	log.Println("Shutting down gracefully...")
	wg.Wait()
	log.Println("Server stopped")
}

func StartGRPCServer(ctx context.Context, repo domain.URLRepository) error {
	RetrieveURLService := application.NewRetrieveURLService(repo)
	gRPCHandler := grpcadapter.NewGRPCServer(*RetrieveURLService)

	server := grpc.NewServer()
	protos.RegisterGetURLServer(server, gRPCHandler)

	port := ":" + os.Getenv("URL_SHORTENER_SVC_GRPC_PORT")

	list, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("cannot listen on tcp connection at %s :%w", port, err)
	}

	log.Printf("gRPC server running on %s", port)

	go func() {
		<-ctx.Done()
		log.Println("Stopping gRPC server...")
		server.GracefulStop()
	}()

	return server.Serve(list)
}

func StartHTTPHandler(ctx context.Context, repo domain.URLRepository, publisher application.URLPublisher) error {
	CreateURLService := application.NewCreateURLAndPublishService(repo, publisher)
	handler := httpAdapter.NewHandler(CreateURLService)

	mux := http.NewServeMux()
	handler.Routes(mux)

	port := ":" + os.Getenv("URL_SHORTENER_SVC_PORT")
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
