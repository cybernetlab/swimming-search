package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/live", probe)
	r.Get("/ready", probe)

	return r
}

func probe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
