package handlers

import (
	"net/http"
	"strings"

	"github.com/mosesniyonk/rwandapi/internal/database"
	"github.com/mosesniyonk/rwandapi/internal/models"
)

// Search searches across all geographic entities.
func Search(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		writeError(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	pattern := "%" + q + "%"
	var results []models.SearchResult

	// Search provinces
	rows, err := database.DB.Query("SELECT id, name FROM provinces WHERE name LIKE ?", pattern)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			var name string
			rows.Scan(&id, &name)
			results = append(results, models.SearchResult{
				Type: "province", ID: id, Name: name, Path: name,
			})
		}
	}

	// Search districts
	rows2, err := database.DB.Query(`
		SELECT d.id, d.name, p.name
		FROM districts d JOIN provinces p ON d.province_id = p.id
		WHERE d.name LIKE ?`, pattern)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var id int
			var name, prov string
			rows2.Scan(&id, &name, &prov)
			results = append(results, models.SearchResult{
				Type: "district", ID: id, Name: name, Path: prov + " > " + name,
			})
		}
	}

	// Search sectors
	rows3, err := database.DB.Query(`
		SELECT s.id, s.name, d.name, p.name
		FROM sectors s
		JOIN districts d ON s.district_id = d.id
		JOIN provinces p ON d.province_id = p.id
		WHERE s.name LIKE ?`, pattern)
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			var id int
			var name, dist, prov string
			rows3.Scan(&id, &name, &dist, &prov)
			results = append(results, models.SearchResult{
				Type: "sector", ID: id, Name: name,
				Path: prov + " > " + dist + " > " + name,
			})
		}
	}

	// Search cells
	rows4, err := database.DB.Query(`
		SELECT c.id, c.name, s.name, d.name, p.name
		FROM cells c
		JOIN sectors s ON c.sector_id = s.id
		JOIN districts d ON s.district_id = d.id
		JOIN provinces p ON d.province_id = p.id
		WHERE c.name LIKE ?`, pattern)
	if err == nil {
		defer rows4.Close()
		for rows4.Next() {
			var id int
			var name, sec, dist, prov string
			rows4.Scan(&id, &name, &sec, &dist, &prov)
			results = append(results, models.SearchResult{
				Type: "cell", ID: id, Name: name,
				Path: prov + " > " + dist + " > " + sec + " > " + name,
			})
		}
	}

	// Search villages
	rows5, err := database.DB.Query(`
		SELECT v.id, v.name, c.name, s.name, d.name, p.name
		FROM villages v
		JOIN cells c ON v.cell_id = c.id
		JOIN sectors s ON c.sector_id = s.id
		JOIN districts d ON s.district_id = d.id
		JOIN provinces p ON d.province_id = p.id
		WHERE v.name LIKE ?`, pattern)
	if err == nil {
		defer rows5.Close()
		for rows5.Next() {
			var id int
			var name, cell, sec, dist, prov string
			rows5.Scan(&id, &name, &cell, &sec, &dist, &prov)
			results = append(results, models.SearchResult{
				Type: "village", ID: id, Name: name,
				Path: prov + " > " + dist + " > " + sec + " > " + cell + " > " + name,
			})
		}
	}

	writeJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    results,
		Count:   len(results),
	})
}
