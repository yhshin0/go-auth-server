package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Register(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
}
