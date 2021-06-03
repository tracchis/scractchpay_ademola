package main

import (
	"github.com/scratchpay_ademola/pkg/clinic"

	"github.com/go-chi/chi"
)

// initRoutes initialize the routing configuration and return a prepared http.Handler
func initRoutes(fetcher clinic.DataFetcher) *chi.Mux {
	mux := chi.NewMux()

	mux.Route("/v1/clinics", func(r chi.Router) {
		r.Post("/search", clinic.Search(fetcher))
		r.Get("/", clinic.GetAllClinics(fetcher))
	})

	return mux
}
