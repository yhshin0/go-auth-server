package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/yhshin0/go-auth-server/internal/config"
)

func Register(r *chi.Mux) {
	cfg := config.GetInstance()

	r.Use(middleware.RequestID)      // 모든 로그에 ID 부여 (가장 먼저)
	r.Use(middleware.RealIP)         // IP 교체 (Logger 전에)
	r.Use(middleware.Recoverer)      // panic 방어 (가능한 바깥쪽)
	r.Use(middleware.Logger)         // 요청 로깅 (위 정보들 활용)
	r.Use(cors.Handler(cors.Options{ // CORS preflight 빠르게 응답
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   cfg.CORS.AllowedMethods,
		AllowCredentials: cfg.CORS.AllowedCredentials,
	}))
	r.Use(middleware.Timeout(cfg.Server.HttpHandlerTimeout)) // 핸들러 타임아웃
}
