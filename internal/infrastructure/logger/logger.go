package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/yhshin0/go-auth-server/internal/defs/enum"
)

var once sync.Once

func Setup(env string) {
	var logger *slog.Logger
	once.Do(func() {
		switch strings.ToLower(env) {
		case enum.ServerEnvLocal.String():
			logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
		case enum.ServerEnvDev.String(), enum.ServerEnvProd.String():
			logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		default:
			err := fmt.Errorf("invalid env type: %s", env)
			panic(err)
		}
	})

	slog.SetDefault(logger)
}
