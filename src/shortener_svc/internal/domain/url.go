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

// URLRepository used to define the accepted behavior of the URL object on a database
type URLRepository interface {
	Persist(url URL) (URL, error)
	GetByShortCode(code string) (URL, error)
}

var (
	ErrCodeEmpty   = errors.New("code cannot be empty")
	ErrRawURLEmpty = errors.New("rawURL cannot be empty")
)
