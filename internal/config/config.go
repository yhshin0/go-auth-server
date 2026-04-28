package config

import (
	"log"
	"sync"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	CORS   struct {
		AllowedOrigins     []string `env:"CORS_ALLOWED_ORIGINS" envSeparator:"," envDefault:"https://example.com:8080"`
		AllowedMethods     []string `env:"CORS_ALLOWED_METHODS" envSeparator:"," envDefault:"GET,HEAD,PUT,PATCH,POST,DELETE,OPTIONS"`
		AllowedCredentials bool     `env:"CORS_ALLOWED_CREDENTIALS" envDefault:"false"`
	}
	DB    DBConfig
	Cache CacheConfig
}

type ServerConfig struct {
	Env                   string        `env:"ENV" envDefault:"local"`
	Port                  string        `env:"PORT" envDefault:"3000"`
	Host                  string        `env:"HOST" envDefault:"localhost"`
	HttpHandlerTimeout    time.Duration `env:"HTTP_HANDLER_TIMEOUT" envDefault:"10s"`
	HttpReadTimeout       time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"5s"`
	HttpReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" envDefault:"2s"`
	HttpWriteTimeout      time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"10s"`
	HttpIdleTimeout       time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
	LogLevel              string        `env:"LOG_LEVEL" envDefault:"info"`
	AccessTokenTTL        time.Duration `env:"ACCESS_TOKEN_TTL" envDefault:"10m"`
	RefreshTokenTTL       time.Duration `env:"REFRESH_TOKEN_TTL" envDefault:"168h"`
	JwtSecret             string        `env:"JWT_SECRET" envDefault:""`
	JwtIssuer             string        `env:"JWT_ISSUER" envDefault:"auth-service"`
	CookieSecure          bool          `env:"COOKIE_SECURE" envDefault:"false"`
	UserSessionLimit      int           `env:"USER_SESSION_LIMIT" envDefault:"5"`
	RefreshLockTTL        time.Duration `env:"REFRESH_LOCK_TTL" envDefault:"3s"` // cache에서 refresh token 갱신 시의 lock time
}

type DBConfig struct {
	Driver       string        `env:"DB_DRIVER" envDefault:"postgres"`
	Host         string        `env:"DB_HOST" envDefault:"localhost"`
	Port         int           `env:"DB_PORT" envDefault:"5432"`
	Name         string        `env:"DB_NAME" envDefault:"my_auth"`
	User         string        `env:"DB_USER" envDefault:"my_user"`
	Password     string        `env:"DB_PASSWORD" envDefault:""`
	SSLMode      string        `env:"DB_SSLMODE" envDefault:"disable"`
	MaxIdleConns int           `env:"DB_MAX_IDLE_CONNS" envDefault:"1"`
	MaxOpenConns int           `env:"DB_MAX_OPEN_CONNS" envDefault:"10"`
	MaxLifetime  time.Duration `env:"DB_MAX_LIFETIME" envDefault:"30m"`
}

type CacheConfig struct {
	Driver    string `env:"CACHE_DRIVER" envDefault:"redis"`
	Host      string `env:"CACHE_HOST" envDefault:"localhost"`
	Port      int    `env:"CACHE_PORT" envDefault:"6379"`
	DB        int    `env:"CACHE_DB" envDefault:"0"`
	Password  string `env:"CACHE_PASSWORD" envDefault:""`
	PoolSize  int    `env:"CACHE_POOL_SIZE" envDefault:"20"`
	KeyPrefix string `env:"CACHE_KEY_PREFIX" envDefault:"auth"`
}

var (
	cfg  Config
	once sync.Once
)

func GetInstance() *Config {
	Setup()
	return &cfg
}

func Setup() {
	once.Do(func() {
		_ = godotenv.Load()

		cfg = Config{}
		if err := env.Parse(&cfg); err != nil {
			// log 설정 전이라서 log 라이브러리 사용
			log.Println("failed to parse config. error:", err.Error())
			panic(err)
		}
	})
}
