package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/mosesniyonk/rwandapi/internal/handlers"
)

// registerRoutes sets up all API routes.
func registerRoutes(r chi.Router) {
	r.Route("/api/v1", func(r chi.Router) {
		// Geographic endpoints
		r.Get("/provinces", handlers.ListProvinces)
		r.Get("/provinces/{id}", handlers.GetProvince)

		r.Get("/districts", handlers.ListDistricts)
		r.Get("/districts/{id}", handlers.GetDistrict)

		r.Get("/sectors", handlers.ListSectors)
		r.Get("/sectors/{id}", handlers.GetSector)

		r.Get("/cells", handlers.ListCells)
		r.Get("/cells/{id}", handlers.GetCell)

		r.Get("/villages", handlers.ListVillages)

		// Banks
		r.Get("/banks", handlers.ListBanks)

		// Utilities
		r.Get("/utilities/electricity/tariffs", handlers.ListElectricityTariffs)
		r.Get("/utilities/water/tariffs", handlers.ListWaterTariffs)

		// Search
		r.Get("/search", handlers.Search)
	})
}
