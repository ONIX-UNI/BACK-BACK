package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func NewPostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
	user := strings.TrimSpace(os.Getenv("POSTGRES_USER"))
	password := os.Getenv("POSTGRES_PASSWORD")
	host := envOrDefault("DB_HOST", envOrDefault("POSTGRES_HOST", "postgres"))
	port := envOrDefault("POSTGRES_PORT", "5432")
	database := envOrDefault("POSTGRES_DB", "postgres")
	sslmode := envOrDefault("POSTGRES_SSLMODE", "disable")

	fmt.Println("USER:", user)
	fmt.Println("PASS:", password)
	fmt.Println("HOST:", host)
	fmt.Println("DB:", database)

	dsnURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   database,
	}

	query := dsnURL.Query()
	query.Set("sslmode", sslmode)
	dsnURL.RawQuery = query.Encode()

	cfg, err := pgxpool.ParseConfig(dsnURL.String())
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	DB = pool
	fmt.Printf("successfully connected to PostgreSQL database: %s\n", database)

	return pool, nil
}

func envOrDefault(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
