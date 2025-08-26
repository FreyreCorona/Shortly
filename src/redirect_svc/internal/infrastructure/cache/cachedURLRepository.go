// Package cache resolves responses for more efficiency and avoid repetitive db calls
package cache

import (
	"errors"

	"github.com/FreyreCorona/Shortly/src/redirect_svc/internal/domain"
)

// URLCacheRepository implement the accepted behavior of Cache systems
type URLCacheRepository interface {
	Get(code string) (*domain.URL, error)
	Set(url domain.URL) error
}

// CachedURLRepository secondary port implements URLRepository and add cache funcionalities
type CachedURLRepository struct {
	Repo  domain.URLRepository
	Cache URLCacheRepository
}

var (
	ErrNoCachedURL = errors.New("specified url not found on cache")
	ErrCacheOff    = errors.New("cache system offline")
)

func NewCachedURLRepository(repo domain.URLRepository, cache URLCacheRepository) *CachedURLRepository {
	return &CachedURLRepository{Repo: repo, Cache: cache}
}

func (c CachedURLRepository) GetByShortCode(code string) (domain.URL, error) {
	cached, err := c.Cache.Get(code)
	if err == nil && cached != nil {
		return *cached, nil
	}
	if err != nil && err != ErrNoCachedURL {
		return domain.URL{}, err
	}

	url, err := c.Repo.GetByShortCode(code)
	if err != nil {
		return domain.URL{}, err
	}

	err = c.Cache.Set(url)
	if err != nil {
		return domain.URL{}, nil
	}
	return url, nil
}

func (c CachedURLRepository) Persist(url domain.URL) (domain.URL, error) {
	url, err := c.Repo.Persist(url)
	if err != nil {
		return domain.URL{}, err
	}

	err = c.Cache.Set(url)
	if err != nil {
		return domain.URL{}, nil
	}

	return url, nil
}
