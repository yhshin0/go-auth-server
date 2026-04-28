package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/yhshin0/go-auth-server/internal/config"
	"github.com/yhshin0/go-auth-server/internal/defs"
	"github.com/yhshin0/go-auth-server/internal/infrastructure/logger"
)

func NewDatabase(cfg *config.DBConfig) (*sqlx.DB, error) {
	// postgres://user:password@localhost:5432/db?sslmode=disable
	uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
	db, err := connectDB(uri)
	if err != nil {
		return nil, err
	}

	logger.Info("success to connect to postgres")

	// 커넥션 풀 설정 (성능 최적화)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)
	return db, nil
}

func connectDB(uri string) (*sqlx.DB, error) {
	const maxTries = 5

	var db *sqlx.DB
	var err error

	for i := 0; i < maxTries; i++ {
		db, err = sqlx.Open(defs.PostgresDriver, uri)
		if err != nil {
			logger.Warn("failed to open postgres", "error", err)
			return nil, err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err = db.PingContext(ctx)
		cancel()

		if err == nil {
			return db, nil
		}

		db.Close()

		wait := time.Duration(1<<i) * time.Second // 1s, 2s, 4s, ...
		logger.Warn("failed to ping postgres", "try", i+1, "wait", wait, "error", err)
		time.Sleep(wait)
	}

	return nil, fmt.Errorf("failed to connect to database(%s). error: %w", defs.PostgresDriver, err)
}
