package shorturl

import (
	"context"
	"errors"
	"fmt"
	"time"
	"urlShortenerService/domain"
	"urlShortenerService/internal/infrastructure/config"

	"github.com/jackc/pgx/v5"
)

var (
	// deleteExpiredStmt is the prepared statement to delete expired slug / url couple from the database
	deleteExpiredStmt string = "DELETE FROM urls WHERE inserted_at < $1;"
	// getStmt is the prepared statement to retrieve a url given a slug from the database
	getStmt string = "SELECT slug, url, inserted_at FROM urls WHERE slug=$1;"
	// setStmt is the prepared statement to insert a slug / url couple into the database
	setStmt string = "INSERT INTO urls (slug, url, inserted_at) VALUES ($1, $2, $3) ON CONFLICT (slug) DO UPDATE SET inserted_at = $3;"
)

// PSQLStore represents a postgres SQL store
type PSQLStore struct {
	conn *pgx.Conn
}

// NewPSQLStore connects to a database and return it inside a PSQLStore
func NewPSQLStore(connConf config.PSQLConnConfig) (*PSQLStore, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connConf.ToConnString())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	store := &PSQLStore{conn: conn}

	err = store.initTables(ctx)
	if err != nil {
		return nil, err
	}

	return store, nil
}

// initTable initializes the PSQL tables
func (s *PSQLStore) initTables(ctx context.Context) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS urls (
		slug TEXT PRIMARY KEY,
		url TEXT NOT NULL,
		inserted_at TIMESTAMP NOT NULL
	);`

	_, err := s.conn.Exec(ctx, createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// DeleteExpired implements the Store interface
func (s *PSQLStore) DeleteExpired(ctx context.Context, duration time.Duration) error {
	cutoff := time.Now().Add(-duration)
	_, err := s.conn.Exec(ctx, deleteExpiredStmt, cutoff)
	return err
}

// Get implements the Store interface
func (s *PSQLStore) Get(ctx context.Context, slug string) (domain.ShortURL, error) {
	var url domain.ShortURL
	err := s.conn.QueryRow(ctx, getStmt, slug).Scan(&url.Slug, &url.URL, &url.InsertedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return url, ErrNotFound
		}
		return url, err
	}
	return url, nil
}

// Set implements the Store interface
func (s *PSQLStore) Set(ctx context.Context, shortURL domain.ShortURL) error {
	if shortURL.InsertedAt.IsZero() {
		shortURL.InsertedAt = time.Now()
	}
	_, err := s.conn.Exec(ctx, setStmt, shortURL.Slug, shortURL.URL, shortURL.InsertedAt.UTC())
	return err
}

// Close closes the database connection
func (s *PSQLStore) Close() error {
	return s.conn.Close(context.Background())
}
