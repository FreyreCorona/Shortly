package main

import (
	"fmt"
	"log"
	"os"

	"github.com/FreyreCorona/Shortly/src/redirect_svc/internal/infrastructure/cache/valkey"
)

func main() {
	c, err := valkey.NewValkeyCache(fmt.Sprintf("%s:%s", os.Getenv("VALKEY_ADDR"), os.Getenv("VALKEY_PORT")),
		os.Getenv("VALKEY_USERNAME"),
		os.Getenv("VALKEY_PASSWORD"))
	if err != nil {
		log.Fatalf("cache connection error :%s", err.Error())
	}
}
