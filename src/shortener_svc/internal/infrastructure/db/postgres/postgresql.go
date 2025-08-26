// Package db Implements persistenci across multiple types of databases
package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/FreyreCorona/Shortly/src/shortener_svc/internal/domain"
	_ "github.com/lib/pq"
)

type PostgresURLRepository struct {
	DB *sql.DB
}

// NewPostgresDB opens a connection pool verify his integrity and return the URLRepository
func NewPostgresDB(dsn string) (*PostgresURLRepository, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = conn.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return &PostgresURLRepository{DB: conn}, nil
}

// Persist implements URLRepository interface and uses transcaction for execute query generated functions
func (r *PostgresURLRepository) Persist(url domain.URL) (domain.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return domain.URL{}, nil
	}
	defer func() {
		err := tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			_ = err // intencionally ignore the error
		}
	}()

	query := New(r.DB).WithTx(tx)

	u, err := query.SaveURL(ctx, SaveURLParams{
		OriginalUrl: url.RawURL,
		ShortCode:   url.ShortCode,
	})
	if err != nil {
		return domain.URL{}, nil
	}
	err = tx.Commit()
	if err != nil {
		return domain.URL{}, nil
	}
	return domain.URL{ID: u.ID, RawURL: url.RawURL, ShortCode: url.ShortCode, CreatedAt: u.CreatedAt.Time}, nil
}

// GetByShortCode implements URLRepository interface and uses transcaction for execute query generated functions
func (r *PostgresURLRepository) GetByShortCode(code string) (domain.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return domain.URL{}, nil
	}
	defer func() {
		err := tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			_ = err // intencionally ignore the error
		}
	}()

	query := New(r.DB).WithTx(tx)

	url, err := query.GetByShortCode(ctx, code)
	if err != nil {
		return domain.URL{}, err
	}
	err = tx.Commit()
	if err != nil {
		return domain.URL{}, nil
	}
	return domain.URL{ID: url.ID, RawURL: url.OriginalUrl, ShortCode: url.ShortCode, CreatedAt: url.CreatedAt.Time}, nil
}
