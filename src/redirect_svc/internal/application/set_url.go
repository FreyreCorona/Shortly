package application

import "github.com/FreyreCorona/Shortly/src/redirect_svc/internal/domain"

type SetURLService struct {
	cache domain.URLCacheRepository
}

func NewSetURLService(cache domain.URLCacheRepository) *SetURLService {
	return &SetURLService{cache: cache}
}

func (s *SetURLService) SetURL(url domain.URL) error {
	return s.cache.Set(url)
}
