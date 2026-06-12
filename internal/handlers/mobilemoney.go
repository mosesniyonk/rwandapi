package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mosesniyonk/rwandapi/internal/database"
	"github.com/mosesniyonk/rwandapi/internal/models"
)

func ListMobileMoney(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(
		"SELECT id, name, provider, ussd_code, agent_ussd, daily_limit, monthly_limit, created_at FROM mobile_money ORDER BY name",
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query mobile money services")
		return
	}
	defer rows.Close()

	var services []models.MobileMoneyService
	for rows.Next() {
		var s models.MobileMoneyService
		if err := rows.Scan(&s.ID, &s.Name, &s.Provider, &s.USSDCode, &s.AgentUSSD, &s.DailyLimit, &s.MonthlyLimit, &s.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan mobile money")
			return
		}

		feeRows, err := database.DB.Query(
			"SELECT id, tx_type, min_amount, max_amount, fee FROM mobile_money_fees WHERE service_id = ? ORDER BY tx_type, CAST(min_amount AS INTEGER)",
			s.ID,
		)
		if err == nil {
			for feeRows.Next() {
				var f models.MobileMoneyFee
				feeRows.Scan(&f.ID, &f.TxType, &f.MinAmount, &f.MaxAmount, &f.Fee)
				s.Fees = append(s.Fees, f)
			}
			feeRows.Close()
		}

		services = append(services, s)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: services, Count: len(services)})
}

func GetMobileMoney(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var s models.MobileMoneyService
	err := database.DB.QueryRow(
		"SELECT id, name, provider, ussd_code, agent_ussd, daily_limit, monthly_limit, created_at FROM mobile_money WHERE id = ?", id,
	).Scan(&s.ID, &s.Name, &s.Provider, &s.USSDCode, &s.AgentUSSD, &s.DailyLimit, &s.MonthlyLimit, &s.CreatedAt)
	if err != nil {
		writeError(w, http.StatusNotFound, "mobile money service not found")
		return
	}

	feeRows, err := database.DB.Query(
		"SELECT id, tx_type, min_amount, max_amount, fee FROM mobile_money_fees WHERE service_id = ? ORDER BY tx_type, CAST(min_amount AS INTEGER)",
		s.ID,
	)
	if err == nil {
		defer feeRows.Close()
		for feeRows.Next() {
			var f models.MobileMoneyFee
			feeRows.Scan(&f.ID, &f.TxType, &f.MinAmount, &f.MaxAmount, &f.Fee)
			s.Fees = append(s.Fees, f)
		}
	}

	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: s})
}
