package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/arslan-atajykov/shortener-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Init(cfg *config.Config) {
	dsn := fmt.Sprintf(

		"postgress://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSql:", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Database not responding:", err)

	}

	DB = pool
	log.Println("Connected to PostgreSQL")

}
