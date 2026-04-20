package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/yhshin0/go-auth-server/internal/config"
)

func Register(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(config.GetInstance().Server.HttpHandlerTimeout))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.GetInstance().CORS.AllowedOrigins,
		AllowedMethods:   config.GetInstance().CORS.AllowedMethods,
		AllowCredentials: config.GetInstance().CORS.AllowedCredentials,
	}))
}
