package domain

import "testing"

func TestRedirectURLStructure(t *testing.T) {
	rawURL := "https://example.com"
	shortCode := "xyz987"

	u := URL{
		RawURL:    rawURL,
		ShortCode: shortCode,
	}

	if u.RawURL != rawURL {
		t.Errorf("expected raw URL %s, got %s", rawURL, u.RawURL)
	}

	if u.ShortCode != shortCode {
		t.Errorf("expected short code %s, got %s", shortCode, u.ShortCode)
	}
}

func TestRedirectErrors(t *testing.T) {
	if ErrCodeEmpty.Error() != "code cannot be empty" {
		t.Errorf("unexpected error message for ErrCodeEmpty")
	}

	if ErrNoURL.Error() != "specified url not exist" {
		t.Errorf("unexpected error message for ErrNoURL")
	}
}
