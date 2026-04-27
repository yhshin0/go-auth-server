package cache

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/yhshin0/go-auth-server/internal/config"
	"github.com/yhshin0/go-auth-server/internal/defs"
	"github.com/yhshin0/go-auth-server/internal/infrastructure/cache/redis"
)

type Client interface {
	CloseWithLog()
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

var (
	cli  Client
	once sync.Once
)

func NewCache(cfg *config.CacheConfig) Client {
	var err error
	once.Do(func() {
		cli, err = setup(cfg)
		if err != nil {
			log.Println("failed to setup cache:", err)
			panic(err)
		}
	})

	return cli
}

func setup(cfg *config.CacheConfig) (Client, error) {
	var err error
	switch strings.ToLower(cfg.Driver) {
	case defs.RedisDriver:
		cli, err = redis.NewRedis(cfg)

	default:
		err = fmt.Errorf("unsupported cache driver: %s", cfg.Driver)
	}

	return cli, err
}
