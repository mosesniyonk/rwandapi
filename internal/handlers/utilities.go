package handlers

import (
	"net/http"

	"github.com/mosesniyonk/rwandapi/internal/database"
	"github.com/mosesniyonk/rwandapi/internal/models"
)

// ListElectricityTariffs returns EUCL electricity tariff bands.
func ListElectricityTariffs(w http.ResponseWriter, r *http.Request) {
	listTariffs(w, "electricity")
}

// ListWaterTariffs returns WASAC water tariff bands.
func ListWaterTariffs(w http.ResponseWriter, r *http.Request) {
	listTariffs(w, "water")
}

func listTariffs(w http.ResponseWriter, utility string) {
	rows, err := database.DB.Query(
		"SELECT id, utility, category, description, min_usage, max_usage, unit, price_per_unit, currency, created_at FROM tariffs WHERE utility = ? ORDER BY category, min_usage",
		utility,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query tariffs")
		return
	}
	defer rows.Close()

	var tariffs []models.Tariff
	for rows.Next() {
		var t models.Tariff
		if err := rows.Scan(&t.ID, &t.Utility, &t.Category, &t.Description,
			&t.MinUsage, &t.MaxUsage, &t.Unit, &t.PricePerUnit, &t.Currency, &t.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan tariff")
			return
		}
		tariffs = append(tariffs, t)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    tariffs,
		Count:   len(tariffs),
	})
}
