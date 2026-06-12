package handlers

import (
	"net/http"

	"github.com/mosesniyonk/rwandapi/internal/database"
	"github.com/mosesniyonk/rwandapi/internal/models"
)

func ListPublicHolidays(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, date, description, is_movable FROM public_holidays ORDER BY date")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query public holidays")
		return
	}
	defer rows.Close()

	var holidays []models.PublicHoliday
	for rows.Next() {
		var h models.PublicHoliday
		var movable int
		if err := rows.Scan(&h.ID, &h.Name, &h.Date, &h.Description, &movable); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan holiday")
			return
		}
		h.IsMovable = movable == 1
		holidays = append(holidays, h)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: holidays, Count: len(holidays)})
}
