package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mosesniyonk/rwandapi/internal/database"
	"github.com/mosesniyonk/rwandapi/internal/models"
)

// TODO: Connect to BNR (National Bank of Rwanda) live data feed to fetch real-time rates.
// BNR publishes daily exchange rates at https://www.bnr.rw/currency/ — consider scraping
// or using their data feed if one becomes available. Currently serves cached/seeded static rates.

const rateCacheDuration = 1 * time.Hour

// ListExchangeRates returns current BNR exchange rates.
func ListExchangeRates(w http.ResponseWriter, r *http.Request) {
	// Check if cached rates are fresh enough
	var fetchedAt string
	err := database.DB.QueryRow("SELECT fetched_at FROM exchange_rates ORDER BY fetched_at DESC LIMIT 1").Scan(&fetchedAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query exchange rates")
		return
	}

	lastUpdated, err := time.Parse("2006-01-02T15:04:05Z", fetchedAt)
	if err != nil {
		// Try alternate format
		lastUpdated, err = time.Parse(time.RFC3339, fetchedAt)
		if err != nil {
			lastUpdated = time.Time{}
		}
	}

	stale := time.Since(lastUpdated) > rateCacheDuration

	// TODO: If stale, attempt to fetch fresh rates from BNR and update the cache.
	// For now we always return cached/seeded rates.
	if stale {
		refreshExchangeRates()
	}

	// Query all cached rates
	rows, err := database.DB.Query("SELECT id, currency_code, currency_name, buying, selling, middle, date, fetched_at FROM exchange_rates ORDER BY currency_code")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query exchange rates")
		return
	}
	defer rows.Close()

	var rates []models.ExchangeRate
	for rows.Next() {
		var rate models.ExchangeRate
		if err := rows.Scan(&rate.ID, &rate.CurrencyCode, &rate.CurrencyName, &rate.Buying, &rate.Selling, &rate.Middle, &rate.Date, &rate.FetchedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan exchange rate")
			return
		}
		rates = append(rates, rate)
	}

	resp := models.ExchangeRateResponse{
		Rates:       rates,
		LastUpdated: fetchedAt,
		Source:      "BNR (National Bank of Rwanda)",
	}

	if stale {
		resp.Warning = "rates may be outdated; live BNR data feed is not yet connected"
	}

	writeJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    resp,
		Count:   len(rates),
	})
}

// refreshExchangeRates attempts to fetch fresh rates from BNR.
// TODO: Implement actual BNR data fetching. For now this is a no-op placeholder
// that updates the fetched_at timestamp so the staleness warning reflects reality.
func refreshExchangeRates() {
	// Placeholder: in a real implementation, this would:
	// 1. Fetch rates from BNR (e.g., scrape https://www.bnr.rw/currency/ or use an XML/JSON feed)
	// 2. Parse the response
	// 3. Upsert rates into the exchange_rates table
	//
	// For now, we just update the fetched_at to avoid repeated refresh attempts.
	now := time.Now().UTC().Format(time.RFC3339)
	today := time.Now().UTC().Format("2006-01-02")

	// Delete old rates and re-insert with updated timestamp (keeps the static data fresh)
	tx, err := database.DB.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	// Read existing rates
	rows, err := tx.Query("SELECT currency_code, currency_name, buying, selling, middle FROM exchange_rates")
	if err != nil {
		return
	}

	type rateRow struct {
		code, name             string
		buying, selling, mid   float64
	}

	var existing []rateRow
	for rows.Next() {
		var r rateRow
		if err := rows.Scan(&r.code, &r.name, &r.buying, &r.selling, &r.mid); err != nil {
			rows.Close()
			return
		}
		existing = append(existing, r)
	}
	rows.Close()

	if _, err := tx.Exec("DELETE FROM exchange_rates"); err != nil {
		return
	}

	for _, r := range existing {
		if _, err := tx.Exec(
			"INSERT INTO exchange_rates (id, currency_code, currency_name, buying, selling, middle, date, fetched_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			uuid.New().String(), r.code, r.name, r.buying, r.selling, r.mid, today, now,
		); err != nil {
			return
		}
	}

	tx.Commit()
}
