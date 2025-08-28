// Package application used to define the default behavior of the application
package application

import (
	"errors"

	"github.com/FreyreCorona/Shortly/src/redirect_svc/internal/domain"
)

type RedirectionService struct {
	cache domain.URLCacheRepository
	repo  domain.URLRepository
}

func NewRedirectionService(cache domain.URLCacheRepository, repo domain.URLRepository) *RedirectionService {
	return &RedirectionService{repo: repo}
}

func (s *RedirectionService) GetURL(code string) (*domain.URL, error) {
	cached, err := s.cache.Get(code)
	if err != nil && !errors.Is(err, domain.ErrNoCachedURL) {
		return nil, err // unespected
	}
	if cached != nil {
		return cached, nil // cache hit
	}

	// search on source
	url, err := s.repo.GetByShortCode(code)
	if err != nil {
		return nil, err
	}

	err = s.cache.Set(url)
	if err != nil {
		return nil, err
	}

	return &url, nil
}
