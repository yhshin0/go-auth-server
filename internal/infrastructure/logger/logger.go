package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/yhshin0/go-auth-server/internal/defs/enum"
)

var (
	logger *slog.Logger
	once   sync.Once
)

func Setup(env string) {
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
}

func Debug(format string, args ...any) {
	logger.Debug(format, args...)
}

func Info(format string, args ...any) {
	logger.Info(format, args...)
}

func Warn(format string, args ...any) {
	logger.Warn(format, args...)
}

func Error(format string, args ...any) {
	logger.Error(format, args...)
}
