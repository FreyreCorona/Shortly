package application

import (
	"github.com/FreyreCorona/Shortly/src/shortener_svc/internal/domain"
)

// CreateURLAndPublish used for create url and publish in a queue
type CreateURLAndPublish struct {
	service   CreateURLService
	publisher URLPublisher
}

type URLPublisher interface {
	Publish(url domain.URL) error
}

func NewCreateURLAndPublishService(repo domain.URLRepository, publisher URLPublisher) *CreateURLAndPublish {
	s := NewCreateURLService(repo)

	return &CreateURLAndPublish{service: s, publisher: publisher}
}

func (s *CreateURLAndPublish) CreateURL(rawURL string) (domain.URL, error) {
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
