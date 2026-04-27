package redis

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/yhshin0/go-auth-server/internal/config"
)

type Client struct {
	keyPrefix string
	cli       *redis.Client
}

func NewRedis(cfg *config.CacheConfig) (*Client, error) {
	const maxRetries = 5

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
		PoolSize: cfg.PoolSize,
	})

	var err error
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = rdb.Ping(ctx).Err()
		cancel()

		if err == nil {
			log.Println("success to connect to redis")
			return &Client{
				keyPrefix: cfg.KeyPrefix,
				cli:       rdb,
			}, nil
		}

		wait := time.Duration(1<<i) * time.Second // 1s, 2s, 4s, ...
		log.Printf("failed to ping (try %d/%d), retrying in %s, error: %s\n", i+1, maxRetries, wait, err.Error())
		time.Sleep(wait)
	}

	return nil, err
}

func (rc *Client) CloseWithLog() {
	if err := rc.cli.Close(); err != nil {
		log.Printf("failed to close redis client: %s\n", err.Error())
	}
}

func (rc *Client) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return rc.cli.Set(ctx, rc.setPrefix(key), value, ttl).Err()
}

func (rc *Client) Get(ctx context.Context, key string) (string, error) {
	return rc.cli.Get(ctx, rc.setPrefix(key)).Result()
}

func (rc *Client) setPrefix(key string) string {
	if !strings.HasPrefix(key+":", rc.keyPrefix) {
		return fmt.Sprintf("%s:%s", rc.keyPrefix, key)
	}
	return key
}
