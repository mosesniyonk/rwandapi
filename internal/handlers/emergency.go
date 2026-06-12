package handlers

import (
	"net/http"

	"github.com/mosesniyonk/rwandapi/internal/database"
	"github.com/mosesniyonk/rwandapi/internal/models"
)

func ListEmergencyNumbers(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	query := "SELECT id, name, number, category, description FROM emergency_numbers"
	var args []interface{}
	if category != "" {
		query += " WHERE category = ?"
		args = append(args, category)
	}
	query += " ORDER BY category, name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query emergency numbers")
		return
	}
	defer rows.Close()

	var numbers []models.EmergencyNumber
	for rows.Next() {
		var n models.EmergencyNumber
		if err := rows.Scan(&n.ID, &n.Name, &n.Number, &n.Category, &n.Description); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan emergency number")
			return
		}
		numbers = append(numbers, n)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: numbers, Count: len(numbers)})
}
