// Package application used to define the default behavior of the application
package application

import (
	"errors"

	"github.com/FreyreCorona/Shortly/src/redirect_svc/internal/domain"
)

type GetURLService struct {
	cache   domain.URLCacheRepository
	repo    domain.URLRepository
	service *SetURL
}

func NewRedirectionService(cache domain.URLCacheRepository, repo domain.URLRepository) *GetURLService {
	s := NewSetURLService(cache)
	return &GetURLService{cache: cache, repo: repo, service: s}
}

func (s *GetURLService) GetURL(code string) (*domain.URL, error) {
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

	err = s.service.SetURL(url)
	if err != nil {
		return nil, err
	}

	return &url, nil
}
