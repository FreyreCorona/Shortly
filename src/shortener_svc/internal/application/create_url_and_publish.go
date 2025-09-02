package application

import (
	"github.com/FreyreCorona/Shortly/src/shortener_svc/internal/domain"
)

// CreateURLAndPublishService used for create url and publish in a queue
type CreateURLAndPublishService struct {
	service   *CreateURLService
	publisher URLPublisher
}

type URLPublisher interface {
	Publish(url domain.URL) error
}

func New(repo domain.URLRepository, publisher URLPublisher) *CreateURLAndPublishService {
	s := NewCreateURLService(repo)

	return &CreateURLAndPublishService{service: s, publisher: publisher}
}

func (s *CreateURLAndPublishService) CreateURL(rawURL string) (domain.URL, error) {
	url, err := s.service.CreateURL(rawURL)
	if err != nil {
		return domain.URL{}, nil
	}
	err = s.publisher.Publish(url)
	if err != nil {
		return domain.URL{}, nil
	}
	return url, nil
}
