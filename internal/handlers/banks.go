package handlers

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mosesniyonk/rwandapi/internal/database"
	"github.com/mosesniyonk/rwandapi/internal/models"
)

func ListBanks(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, short_name, swift_code, phone, created_at FROM banks ORDER BY name")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query banks")
		return
	}
	defer rows.Close()

	var banks []models.Bank
	for rows.Next() {
		var b models.Bank
		if err := rows.Scan(&b.ID, &b.Name, &b.ShortName, &b.SwiftCode, &b.Phone, &b.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan bank")
			return
		}

		brRows, err := database.DB.Query("SELECT id, bank_id, name, location, phone FROM branches WHERE bank_id = ? ORDER BY name", b.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to query branches")
			return
		}

		for brRows.Next() {
			var br models.Branch
			if err := brRows.Scan(&br.ID, &br.BankID, &br.Name, &br.Location, &br.Phone); err != nil {
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

func GetBank(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var b models.Bank
	err := database.DB.QueryRow("SELECT id, name, short_name, swift_code, phone, created_at FROM banks WHERE id = ?", id).
		Scan(&b.ID, &b.Name, &b.ShortName, &b.SwiftCode, &b.Phone, &b.CreatedAt)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "bank not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query bank")
		return
	}

	rows, err := database.DB.Query("SELECT id, bank_id, name, location, phone FROM branches WHERE bank_id = ? ORDER BY name", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query branches")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var br models.Branch
		if err := rows.Scan(&br.ID, &br.BankID, &br.Name, &br.Location, &br.Phone); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan branch")
			return
		}
		b.Branches = append(b.Branches, br)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: b})
}
