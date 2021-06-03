package main

import (
	"github.com/scratchpay_ademola/pkg/clinic"

	"github.com/go-chi/chi"
)

// initRoutes initialize the routing configuration and return a prepared http.Handler
func initRoutes() *chi.Mux {
	mux := chi.NewMux()

	mux.Route("/v1/clinics", func(r chi.Router) {
		r.Post("/", clinic.Search())
		r.Get("/", clinic.GetAllClinics())
	})

	return mux
}
