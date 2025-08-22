// Package valkey implements interface CacheURLRepository
package valkey

import (
	"context"
	"encoding/json"
	"time"

	"github.com/FreyreCorona/Shortly/url_shortener_svc/internal/domain"
	"github.com/FreyreCorona/Shortly/url_shortener_svc/internal/infrastructure/cache"
	"github.com/valkey-io/valkey-go"
)

type ValkeyCachedURLRepository struct {
	Client valkey.Client
}

func NewValkeyCache(address, username, password string) (*ValkeyCachedURLRepository, error) {
	o := valkey.ClientOption{
		InitAddress:       []string{address},
		Username:          username,
		Password:          password,
		ForceSingleClient: true,
	}
	c, err := valkey.NewClient(o)
	if err != nil {
		return nil, err
	}

	// ping to service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = c.Do(ctx, c.B().Ping().Build()).Error()
	if err != nil {
		return nil, err
	}

	return &ValkeyCachedURLRepository{Client: c}, nil
}

func (c ValkeyCachedURLRepository) Get(code string) (*domain.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := c.Client.Do(ctx, c.Client.B().Exists().Key(code).Build()).ToBool()
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, cache.ErrNoCachedURL
	}

	data, err := c.Client.Do(ctx, c.Client.B().Get().Key(code).Build()).ToString()
	if err != nil {
		return nil, err
	}
	var url domain.URL
	err = json.Unmarshal([]byte(data), &url)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (c ValkeyCachedURLRepository) Set(url domain.URL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	value, err := json.Marshal(url)
	if err != nil {
		return err
	}

	return c.Client.Do(ctx, c.Client.B().Set().Key(url.ShortCode).Value(string(value)).Ex(30*time.Minute).Build()).Error()
}
