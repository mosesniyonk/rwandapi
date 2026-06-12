package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mosesniyonk/rwandapi/internal/database"
	"github.com/mosesniyonk/rwandapi/internal/models"
)

func ListTelecoms(w http.ResponseWriter, r *http.Request) {
	typeFilter := r.URL.Query().Get("type")

	query := "SELECT id, name, short_name, type, website, customer_care, ussd_codes, created_at FROM telecoms"
	var args []interface{}
	if typeFilter != "" {
		query += " WHERE type = ?"
		args = append(args, typeFilter)
	}
	query += " ORDER BY name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query telecoms")
		return
	}
	defer rows.Close()

	var telecoms []models.Telecom
	for rows.Next() {
		var t models.Telecom
		var ussdStr string
		if err := rows.Scan(&t.ID, &t.Name, &t.ShortName, &t.Type, &t.Website, &t.CustomerCare, &ussdStr, &t.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan telecom")
			return
		}

		if ussdStr != "" {
			for _, entry := range strings.Split(ussdStr, ",") {
				parts := strings.SplitN(entry, ":", 2)
				if len(parts) == 2 {
					t.USSDCodes = append(t.USSDCodes, models.USSDCode{Code: parts[0], Description: parts[1]})
				}
			}
		}

		planRows, err := database.DB.Query(
			"SELECT id, name, type, speed, data_cap, price, validity FROM telecom_plans WHERE telecom_id = ? ORDER BY name",
			t.ID,
		)
		if err == nil {
			for planRows.Next() {
				var p models.TelecomPlan
				planRows.Scan(&p.ID, &p.Name, &p.Type, &p.Speed, &p.DataCap, &p.Price, &p.Validity)
				t.Plans = append(t.Plans, p)
			}
			planRows.Close()
		}

		telecoms = append(telecoms, t)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: telecoms, Count: len(telecoms)})
}

func GetTelecom(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var t models.Telecom
	var ussdStr string
	err := database.DB.QueryRow(
		"SELECT id, name, short_name, type, website, customer_care, ussd_codes, created_at FROM telecoms WHERE id = ?", id,
	).Scan(&t.ID, &t.Name, &t.ShortName, &t.Type, &t.Website, &t.CustomerCare, &ussdStr, &t.CreatedAt)
	if err != nil {
		writeError(w, http.StatusNotFound, "telecom not found")
		return
	}

	if ussdStr != "" {
		for _, entry := range strings.Split(ussdStr, ",") {
			parts := strings.SplitN(entry, ":", 2)
			if len(parts) == 2 {
				t.USSDCodes = append(t.USSDCodes, models.USSDCode{Code: parts[0], Description: parts[1]})
			}
		}
	}

	planRows, err := database.DB.Query(
		"SELECT id, name, type, speed, data_cap, price, validity FROM telecom_plans WHERE telecom_id = ? ORDER BY name",
		t.ID,
	)
	if err == nil {
		defer planRows.Close()
		for planRows.Next() {
			var p models.TelecomPlan
			planRows.Scan(&p.ID, &p.Name, &p.Type, &p.Speed, &p.DataCap, &p.Price, &p.Validity)
			t.Plans = append(t.Plans, p)
		}
	}

	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: t})
}
