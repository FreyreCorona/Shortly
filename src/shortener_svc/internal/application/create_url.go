// Package application use cases for the URL entities
package application

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/FreyreCorona/Shortly/src/shortener_svc/internal/domain"
)

// CreateURLService use case service for generate short code by URL
type CreateURLService struct {
	repo domain.URLRepository
}

func NewCreateURLService(repo domain.URLRepository) *CreateURLService {
	return &CreateURLService{repo: repo}
}

// CreateURL calls the repository for persist the current object URL created be rawURL
func (s CreateURLService) CreateURL(rawURL string) (domain.URL, error) {
	if rawURL == "" {
		return domain.URL{}, domain.ErrRawURLEmpty
	}

	shortCode, err := generateCode(6)
	if err != nil {
		return domain.URL{}, err
	}
	url := domain.URL{RawURL: rawURL, ShortCode: shortCode}
	url, err = s.repo.Persist(url)
	if err != nil {
		return domain.URL{}, nil
	}
	return url, nil
}

// generateCode creates an array of 6 bytes and fills with random bytes and make url-safe
func generateCode(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	code := base64.RawURLEncoding.EncodeToString(b)
	return code, nil
}
