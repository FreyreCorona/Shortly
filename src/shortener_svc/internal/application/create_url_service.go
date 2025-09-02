package application

import "github.com/FreyreCorona/Shortly/src/shortener_svc/internal/domain"

type CreateURLService interface {
	CreateURL(rawURL string) (domain.URL, error)
}
