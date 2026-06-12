package handlers

import (
	"net/http"

	"github.com/mosesniyonk/rwandapi/internal/database"
	"github.com/mosesniyonk/rwandapi/internal/models"
)

// ListBanks returns all banks with their branches.
func ListBanks(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, short_name, swift_code, created_at FROM banks ORDER BY name")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query banks")
		return
	}
	defer rows.Close()

	var banks []models.Bank
	for rows.Next() {
		var b models.Bank
		if err := rows.Scan(&b.ID, &b.Name, &b.ShortName, &b.SwiftCode, &b.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan bank")
			return
		}

		// Fetch branches for this bank
		brRows, err := database.DB.Query("SELECT id, bank_id, name, location FROM branches WHERE bank_id = ? ORDER BY name", b.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to query branches")
			return
		}

		for brRows.Next() {
			var br models.Branch
			if err := brRows.Scan(&br.ID, &br.BankID, &br.Name, &br.Location); err != nil {
				brRows.Close()
				writeError(w, http.StatusInternalServerError, "failed to scan branch")
				return
			}
			b.Branches = append(b.Branches, br)
		}
		brRows.Close()

		banks = append(banks, b)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    banks,
		Count:   len(banks),
	})
}
