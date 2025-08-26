package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/FreyreCorona/Shortly/shortener_write/internal/application"
	"github.com/FreyreCorona/Shortly/shortener_write/internal/infrastructure/cache"
	"github.com/FreyreCorona/Shortly/shortener_write/internal/infrastructure/cache/valkey"
	db "github.com/FreyreCorona/Shortly/shortener_write/internal/infrastructure/db/postgres"
	httpAddapter "github.com/FreyreCorona/Shortly/shortener_write/internal/infrastructure/http"
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
	c, err := valkey.NewValkeyCache(fmt.Sprintf("%s:%s", os.Getenv("VALKEY_ADDR"), os.Getenv("VALKEY_PORT")),
		os.Getenv("VALKEY_USERNAME"),
		os.Getenv("VALKEY_PASSWORD"))
	if err != nil {
		log.Fatalf("cache connection error :%s", err.Error())
	}

	r := cache.NewCachedURLRepository(repo, c)
	service := application.NewURLService(r)
	handler := httpAddapter.NewHandler(service)

	mux := http.NewServeMux()
	handler.Routes(mux)
	runningPort, err := strconv.Atoi(os.Getenv("URL_SHORTENER_SVC_PORT"))
	if err != nil {
		log.Fatal("Uknown port")
	}
	fmt.Printf("Server running on : %d \n", runningPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", runningPort), mux))
}
