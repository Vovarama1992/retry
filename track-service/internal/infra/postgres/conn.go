package postgres

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/Vovarama1992/go-utils/pgutil"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/sony/gobreaker"
)

func NewPgConn() *sql.DB {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL не задан")
	}

	pool, err := pgutil.NewPool(dbURL, pgutil.DBPoolConfig{
		MaxConnLifetime:   30 * time.Minute,
		MaxConnIdleTime:   10 * time.Minute,
		HealthCheckPeriod: 2 * time.Minute,
		ConnectTimeout:    5 * time.Second,
		MaxConns:          10,
	})
	if err != nil {
		log.Fatalf("ошибка подключения к БД: %v", err)
	}

	db := stdlib.OpenDBFromPool(pool)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("sql.DB ping failed: %v", err)
	}

	return db
}

func NewPgBreaker() *gobreaker.CircuitBreaker {
	return pgutil.NewBreaker(pgutil.BreakerConfig{
		Name:             "track-service", // фикс под этот сервис
		OpenTimeout:      10 * time.Second,
		FailureThreshold: 5,
		MaxRequests:      3,
	})
}
