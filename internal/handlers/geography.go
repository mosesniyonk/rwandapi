package handlers

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mosesniyonk/rwandapi/internal/database"
	"github.com/mosesniyonk/rwandapi/internal/models"
)

func ListProvinces(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, created_at FROM provinces ORDER BY name")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query provinces")
		return
	}
	defer rows.Close()

	var provinces []models.Province
	for rows.Next() {
		var p models.Province
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan province")
			return
		}
		provinces = append(provinces, p)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    provinces,
		Count:   len(provinces),
	})
}

func GetProvince(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var p models.Province
	err := database.DB.QueryRow("SELECT id, name, created_at FROM provinces WHERE id = ?", id).
		Scan(&p.ID, &p.Name, &p.CreatedAt)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "province not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query province")
		return
	}

	rows, err := database.DB.Query("SELECT id, province_id, name, created_at FROM districts WHERE province_id = ? ORDER BY name", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query districts")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var d models.District
		if err := rows.Scan(&d.ID, &d.ProvinceID, &d.Name, &d.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan district")
			return
		}
		p.Districts = append(p.Districts, d)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: p})
}

func ListDistricts(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, province_id, name, created_at FROM districts"
	var args []interface{}

	if pid := r.URL.Query().Get("province_id"); pid != "" {
		query += " WHERE province_id = ?"
		args = append(args, pid)
	}
	query += " ORDER BY name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query districts")
		return
	}
	defer rows.Close()

	var districts []models.District
	for rows.Next() {
		var d models.District
		if err := rows.Scan(&d.ID, &d.ProvinceID, &d.Name, &d.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan district")
			return
		}
		districts = append(districts, d)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    districts,
		Count:   len(districts),
	})
}

func GetDistrict(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var d models.District
	err := database.DB.QueryRow("SELECT id, province_id, name, created_at FROM districts WHERE id = ?", id).
		Scan(&d.ID, &d.ProvinceID, &d.Name, &d.CreatedAt)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "district not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query district")
		return
	}

	rows, err := database.DB.Query("SELECT id, district_id, name, created_at FROM sectors WHERE district_id = ? ORDER BY name", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query sectors")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var s models.Sector
		if err := rows.Scan(&s.ID, &s.DistrictID, &s.Name, &s.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan sector")
			return
		}
		d.Sectors = append(d.Sectors, s)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: d})
}

func ListSectors(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, district_id, name, created_at FROM sectors"
	var args []interface{}

	if did := r.URL.Query().Get("district_id"); did != "" {
		query += " WHERE district_id = ?"
		args = append(args, did)
	}
	query += " ORDER BY name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query sectors")
		return
	}
	defer rows.Close()

	var sectors []models.Sector
	for rows.Next() {
		var s models.Sector
		if err := rows.Scan(&s.ID, &s.DistrictID, &s.Name, &s.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan sector")
			return
		}
		sectors = append(sectors, s)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    sectors,
		Count:   len(sectors),
	})
}

func GetSector(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var s models.Sector
	err := database.DB.QueryRow("SELECT id, district_id, name, created_at FROM sectors WHERE id = ?", id).
		Scan(&s.ID, &s.DistrictID, &s.Name, &s.CreatedAt)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "sector not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query sector")
		return
	}

	rows, err := database.DB.Query("SELECT id, sector_id, name, created_at FROM cells WHERE sector_id = ? ORDER BY name", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query cells")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Cell
		if err := rows.Scan(&c.ID, &c.SectorID, &c.Name, &c.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan cell")
			return
		}
		s.Cells = append(s.Cells, c)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: s})
}

func ListCells(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, sector_id, name, created_at FROM cells"
	var args []interface{}

	if sid := r.URL.Query().Get("sector_id"); sid != "" {
		query += " WHERE sector_id = ?"
		args = append(args, sid)
	}
	query += " ORDER BY name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query cells")
		return
	}
	defer rows.Close()

	var cells []models.Cell
	for rows.Next() {
		var c models.Cell
		if err := rows.Scan(&c.ID, &c.SectorID, &c.Name, &c.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan cell")
			return
		}
		cells = append(cells, c)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    cells,
		Count:   len(cells),
	})
}

func GetCell(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var c models.Cell
	err := database.DB.QueryRow("SELECT id, sector_id, name, created_at FROM cells WHERE id = ?", id).
		Scan(&c.ID, &c.SectorID, &c.Name, &c.CreatedAt)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "cell not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query cell")
		return
	}

	rows, err := database.DB.Query("SELECT id, cell_id, name, created_at FROM villages WHERE cell_id = ? ORDER BY name", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query villages")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var v models.Village
		if err := rows.Scan(&v.ID, &v.CellID, &v.Name, &v.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan village")
			return
		}
		c.Villages = append(c.Villages, v)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: c})
}

func ListVillages(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, cell_id, name, created_at FROM villages"
	var args []interface{}

	if cid := r.URL.Query().Get("cell_id"); cid != "" {
		query += " WHERE cell_id = ?"
		args = append(args, cid)
	}
	query += " ORDER BY name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query villages")
		return
	}
	defer rows.Close()

	var villages []models.Village
	for rows.Next() {
		var v models.Village
		if err := rows.Scan(&v.ID, &v.CellID, &v.Name, &v.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan village")
			return
		}
		villages = append(villages, v)
	}

	writeJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    villages,
		Count:   len(villages),
	})
}
