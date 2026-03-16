package domain

import (
	"testing"
	"time"
)

func TestURLStructure(t *testing.T) {
	id := int64(1)
	rawURL := "https://shortly.io"
	shortCode := "abc123"
	now := time.Now()

	u := URL{
		ID:        id,
		RawURL:    rawURL,
		ShortCode: shortCode,
		CreatedAt: now,
	}

	if u.ID != id {
		t.Errorf("expected ID %d, got %d", id, u.ID)
	}

	if u.RawURL != rawURL {
		t.Errorf("expected RawURL %s, got %s", rawURL, u.RawURL)
	}

	if u.ShortCode != shortCode {
		t.Errorf("expected ShortCode %s, got %s", shortCode, u.ShortCode)
	}

	if !u.CreatedAt.Equal(now) {
		t.Errorf("expected CreatedAt %v, got %v", now, u.CreatedAt)
	}
}

func TestDomainErrors(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		expected string
	}{
		{"ErrCodeEmpty", ErrCodeEmpty, "code cannot be empty"},
		{"ErrRawURLEmpty", ErrRawURLEmpty, "rawURL cannot be empty"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err.Error() != tc.expected {
				t.Errorf("%s: expected error message %q, got %q", tc.name, tc.expected, tc.err.Error())
			}
		})
	}
}
