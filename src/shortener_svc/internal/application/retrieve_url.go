package application

import "github.com/FreyreCorona/Shortly/src/shortener_svc/internal/domain"

// RetrieveURLService use case for retrieving the url object from the database
type RetrieveURLService struct {
	repo domain.URLRepository
}

func NewRetrieveURLService(repo domain.URLRepository) *RetrieveURLService {
	return &RetrieveURLService{repo: repo}
}

// GetURL calls the repository for retrieve the URL object
func (s RetrieveURLService) GetURL(code string) (domain.URL, error) {
	if code == "" {
		return domain.URL{}, domain.ErrCodeEmpty
	}

	url, err := s.repo.GetByShortCode(code)
	if err != nil {
		return domain.URL{}, err
	}

	return url, nil
}
