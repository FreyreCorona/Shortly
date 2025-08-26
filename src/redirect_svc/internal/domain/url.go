// Package domain represent the pure bussines rules and entities
package domain

import (
	"errors"
	"time"
)

// URL the main type of the Shortly
type URL struct {
	ID        int64
	RawURL    string
	ShortCode string
	CreatedAt time.Time
}

// URLCacheRepository used to define the accepted behavior of the URL object on cache
type URLCacheRepository interface {
	Get(code string) (*URL, error)
	Set(url URL) error
}

// URLRepository used to define the accepted behavior of URL source
type URLRepository interface {
	GetByShortCode(code string) (URL, error)
}

var (
	ErrCodeEmpty   = errors.New("code cannot be empty")
	ErrNoCachedURL = errors.New("specified url not found on cache")
	ErrNoURL       = errors.New("specified url not exist")
)
