package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Router(r *chi.Mux /*, m *metrics.HTTPServer*/) {
	// v1 := ver1.New(uc)

	r.Handle("/metrics", promhttp.Handler())

	// r.Route("/mnepryakhin/my-app/api", func(r chi.Router) {
	// 	r.Use(logger.Middleware)
	// 	r.Use(metrics.NewMiddleware(m))
	// 	r.Use(otel.Middleware)

	// 	r.Route("/v1", func(r chi.Router) {
	// 		r.Post("/profile", v1.CreateProfile)
	// 		r.Put("/profile", v1.UpdateProfile)
	// 		r.Get("/profile/{id}", v1.GetProfile)
	// 		r.Delete("/profile/{id}", v1.DeleteProfile)
	// 	})
	// })
}
