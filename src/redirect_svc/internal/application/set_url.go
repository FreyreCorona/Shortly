package application

import "github.com/FreyreCorona/Shortly/src/redirect_svc/internal/domain"

type SetURL struct {
	cache domain.URLCacheRepository
}

func NewSetURLService(cache domain.URLCacheRepository) *SetURL {
	return &SetURL{cache: cache}
}

func (s *SetURL) SetURL(url domain.URL) error {
	return s.cache.Set(url)
}
