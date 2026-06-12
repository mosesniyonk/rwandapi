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
		r.Get("/banks/{id}", handlers.GetBank)

		// Exchange rates
		r.Get("/exchange-rates", handlers.ListExchangeRates)

		// Utilities
		r.Get("/utilities/electricity/tariffs", handlers.ListElectricityTariffs)
		r.Get("/utilities/water/tariffs", handlers.ListWaterTariffs)

		// Telecoms & ISPs
		r.Get("/telecoms", handlers.ListTelecoms)
		r.Get("/telecoms/{id}", handlers.GetTelecom)

		// Mobile Money
		r.Get("/mobile-money", handlers.ListMobileMoney)
		r.Get("/mobile-money/{id}", handlers.GetMobileMoney)

		// Emergency & Useful Numbers
		r.Get("/emergency-numbers", handlers.ListEmergencyNumbers)

		// Public Holidays
		r.Get("/holidays", handlers.ListPublicHolidays)

		// Country Info
		r.Get("/info", handlers.GetCountryInfo)

		// Search
		r.Get("/search", handlers.Search)
	})
}
