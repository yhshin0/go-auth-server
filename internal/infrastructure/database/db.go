package database

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/yhshin0/go-auth-server/internal/config"
	"github.com/yhshin0/go-auth-server/internal/defs"
	"github.com/yhshin0/go-auth-server/internal/infrastructure/database/postgres"
	"github.com/yhshin0/go-auth-server/internal/infrastructure/logger"
)

type DB struct {
	*sqlx.DB
}

func NewDatabase(cfg *config.DBConfig) (*DB, error) {
	switch strings.ToLower(cfg.Driver) {
	case defs.PostgresDriver:
		db, err := postgres.NewDatabase(cfg)
		if err != nil {
			return nil, err
		}
		return &DB{db}, nil
	default:
		err := fmt.Errorf("unsupported database driver: %s", cfg.Driver)
		return nil, err
	}
}

func (db *DB) CloseWithLog() {
	if err := db.DB.Close(); err != nil {
		logger.Warn("failed to close database", "error", err)
	}
}
