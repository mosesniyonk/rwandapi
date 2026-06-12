package database

import (
	"database/sql"
	"log"
)

// Seed populates the database with initial data if tables are empty.
func Seed(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM provinces").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		log.Println("Database already seeded, skipping.")
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := seedProvinces(tx); err != nil {
		return err
	}
	if err := seedDistricts(tx); err != nil {
		return err
	}
	if err := seedSectors(tx); err != nil {
		return err
	}
	if err := seedCells(tx); err != nil {
		return err
	}
	if err := seedVillages(tx); err != nil {
		return err
	}
	if err := seedBanks(tx); err != nil {
		return err
	}
	if err := seedTariffs(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func seedProvinces(tx *sql.Tx) error {
	provinces := []string{
		"Kigali City",
		"Eastern Province",
		"Northern Province",
		"Southern Province",
		"Western Province",
	}
	for _, name := range provinces {
		if _, err := tx.Exec("INSERT INTO provinces (name) VALUES (?)", name); err != nil {
			return err
		}
	}
	return nil
}

func seedDistricts(tx *sql.Tx) error {
	// province name -> list of districts
	districts := map[string][]string{
		"Kigali City": {"Gasabo", "Kicukiro", "Nyarugenge"},
		"Eastern Province": {
			"Bugesera", "Gatsibo", "Kayonza", "Kirehe", "Ngoma",
			"Nyagatare", "Rwamagana",
		},
		"Northern Province": {
			"Burera", "Gakenke", "Gicumbi", "Musanze", "Rulindo",
		},
		"Southern Province": {
			"Gisagara", "Huye", "Kamonyi", "Muhanga", "Nyamagabe",
			"Nyanza", "Nyaruguru", "Ruhango",
		},
		"Western Province": {
			"Karongi", "Ngororero", "Nyabihu", "Nyamasheke",
			"Rubavu", "Rusizi", "Rutsiro",
		},
	}

	for province, dists := range districts {
		var pid int
		if err := tx.QueryRow("SELECT id FROM provinces WHERE name = ?", province).Scan(&pid); err != nil {
			return err
		}
		for _, d := range dists {
			if _, err := tx.Exec("INSERT INTO districts (province_id, name) VALUES (?, ?)", pid, d); err != nil {
				return err
			}
		}
	}
	return nil
}

func seedSectors(tx *sql.Tx) error {
	// Seed sectors for Gasabo and Musanze districts (representative sample)
	sectors := map[string][]string{
		"Gasabo": {
			"Bumbogo", "Gatsata", "Gikomero", "Gisozi", "Jabana",
			"Jali", "Kacyiru", "Kimihurura", "Kimironko", "Kinyinya",
			"Ndera", "Nduba", "Remera", "Rusororo", "Rutunga",
		},
		"Musanze": {
			"Busogo", "Cyuve", "Gacaca", "Gashaki", "Gataraga",
			"Kimonyi", "Kinigi", "Muhoza", "Muko", "Musanze",
			"Nkotsi", "Nyange", "Remera", "Rwaza", "Shingiro",
		},
	}

	for district, secs := range sectors {
		var did int
		if err := tx.QueryRow("SELECT id FROM districts WHERE name = ?", district).Scan(&did); err != nil {
			return err
		}
		for _, s := range secs {
			if _, err := tx.Exec("INSERT INTO sectors (district_id, name) VALUES (?, ?)", did, s); err != nil {
				return err
			}
		}
	}
	return nil
}

func seedCells(tx *sql.Tx) error {
	// Seed cells for a few Gasabo sectors
	cells := map[string][]string{
		"Kimironko": {"Bibare", "Kibagabaga", "Kimironko", "Nyagatovu"},
		"Remera":    {"Gishushu", "Nyabisindu", "Rukiri I", "Rukiri II"},
		"Kacyiru":   {"Kamatamu", "Kamutwa", "Kibaza", "Rugando"},
		// A few from Musanze
		"Muhoza":  {"Cyivugiza", "Kampanga", "Kimonyi", "Mpenge", "Ruhongore"},
		"Kimonyi": {"Gasiza", "Kabeza", "Kavumu", "Nyabageni"},
	}

	for sector, cls := range cells {
		var sid int
		if err := tx.QueryRow("SELECT id FROM sectors WHERE name = ? LIMIT 1", sector).Scan(&sid); err != nil {
			// skip if sector not found (e.g. duplicate name resolution)
			continue
		}
		for _, c := range cls {
			if _, err := tx.Exec("INSERT INTO cells (sector_id, name) VALUES (?, ?)", sid, c); err != nil {
				return err
			}
		}
	}
	return nil
}

func seedVillages(tx *sql.Tx) error {
	// Seed villages for a few cells in Kimironko
	villages := map[string][]string{
		"Kibagabaga":  {"Ingenzi", "Umucyo", "Urumuri", "Ikaze"},
		"Kimironko":   {"Ituze", "Amahoro", "Umwezi", "Imboni"},
		"Bibare":      {"Intsinzi", "Izuba", "Urukundo"},
		"Nyagatovu":   {"Imena", "Isibo", "Ubumwe"},
	}

	for cell, vils := range villages {
		var cid int
		if err := tx.QueryRow("SELECT id FROM cells WHERE name = ? LIMIT 1", cell).Scan(&cid); err != nil {
			continue
		}
		for _, v := range vils {
			if _, err := tx.Exec("INSERT INTO villages (cell_id, name) VALUES (?, ?)", cid, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func seedBanks(tx *sql.Tx) error {
	type bankData struct {
		Name      string
		ShortName string
		SwiftCode string
		Branches  []struct{ Name, Location string }
	}

	banks := []bankData{
		{
			Name: "Bank of Kigali", ShortName: "BK", SwiftCode: "BKIGRWRW",
			Branches: []struct{ Name, Location string }{
				{"BK Headquarters", "Kigali, KN 4 Ave"},
				{"BK Remera", "Kigali, Remera"},
				{"BK Nyabugogo", "Kigali, Nyabugogo"},
				{"BK Musanze", "Musanze"},
				{"BK Huye", "Huye"},
			},
		},
		{
			Name: "Banque Populaire du Rwanda", ShortName: "BPR", SwiftCode: "BPRBRWRW",
			Branches: []struct{ Name, Location string }{
				{"BPR Headquarters", "Kigali, KN 67 St"},
				{"BPR Kicukiro", "Kigali, Kicukiro"},
				{"BPR Rubavu", "Rubavu"},
			},
		},
		{
			Name: "Equity Bank Rwanda", ShortName: "Equity", SwiftCode: "EABORWRW",
			Branches: []struct{ Name, Location string }{
				{"Equity Headquarters", "Kigali, KG 644 St"},
				{"Equity Remera", "Kigali, Remera"},
			},
		},
		{
			Name: "I&M Bank Rwanda", ShortName: "I&M", SwiftCode: "IMABORWRW",
			Branches: []struct{ Name, Location string }{
				{"I&M Headquarters", "Kigali, KN 3 Ave"},
				{"I&M Kacyiru", "Kigali, Kacyiru"},
			},
		},
		{
			Name: "KCB Bank Rwanda", ShortName: "KCB", SwiftCode: "ABORWRW1",
			Branches: []struct{ Name, Location string }{
				{"KCB Headquarters", "Kigali, KN 63 St"},
				{"KCB Nyarugenge", "Kigali, Nyarugenge"},
			},
		},
		{
			Name: "COGEBANQUE", ShortName: "COGEBANQUE", SwiftCode: "COGERWRW",
			Branches: []struct{ Name, Location string }{
				{"COGEBANQUE Headquarters", "Kigali, KN 2 Ave"},
			},
		},
		{
			Name: "Guaranty Trust Bank Rwanda", ShortName: "GTBank", SwiftCode: "GTBIRWRW",
			Branches: []struct{ Name, Location string }{
				{"GTBank Headquarters", "Kigali, KG 5 Ave"},
			},
		},
		{
			Name: "Access Bank Rwanda", ShortName: "Access", SwiftCode: "ABORWRWX",
			Branches: []struct{ Name, Location string }{
				{"Access Bank Headquarters", "Kigali, KG 8 Ave"},
			},
		},
		{
			Name: "Ecobank Rwanda", ShortName: "Ecobank", SwiftCode: "ECABORWRW",
			Branches: []struct{ Name, Location string }{
				{"Ecobank Headquarters", "Kigali, KN 4 Ave"},
			},
		},
		{
			Name: "National Bank of Rwanda", ShortName: "BNR", SwiftCode: "ABORWRWR",
			Branches: []struct{ Name, Location string }{
				{"BNR Headquarters", "Kigali, KN 6 Ave"},
			},
		},
	}

	for _, b := range banks {
		res, err := tx.Exec("INSERT INTO banks (name, short_name, swift_code) VALUES (?, ?, ?)",
			b.Name, b.ShortName, b.SwiftCode)
		if err != nil {
			return err
		}
		bankID, _ := res.LastInsertId()
		for _, br := range b.Branches {
			if _, err := tx.Exec("INSERT INTO branches (bank_id, name, location) VALUES (?, ?, ?)",
				bankID, br.Name, br.Location); err != nil {
				return err
			}
		}
	}
	return nil
}

func seedTariffs(tx *sql.Tx) error {
	type tariffData struct {
		Utility      string
		Category     string
		Description  string
		MinUsage     float64
		MaxUsage     float64
		Unit         string
		PricePerUnit float64
	}

	tariffs := []tariffData{
		// EUCL Electricity tariffs (RWF/kWh)
		{"electricity", "residential", "Residential: 0-15 kWh (lifeline)", 0, 15, "kWh", 89},
		{"electricity", "residential", "Residential: 15-50 kWh", 15, 50, "kWh", 212},
		{"electricity", "residential", "Residential: 50-100 kWh", 50, 100, "kWh", 255},
		{"electricity", "residential", "Residential: >100 kWh", 100, 0, "kWh", 277},
		{"electricity", "commercial", "Non-residential / Commercial", 0, 0, "kWh", 261},
		{"electricity", "industrial", "Industrial: Medium Voltage", 0, 0, "kWh", 189},
		{"electricity", "industrial", "Industrial: High Voltage", 0, 0, "kWh", 133},
		{"electricity", "telecom_towers", "Telecom Towers", 0, 0, "kWh", 261},

		// WASAC Water tariffs (RWF/m3)
		{"water", "residential", "Public standpipe", 0, 0, "m3", 293},
		{"water", "residential", "Residential: 0-5 m3", 0, 5, "m3", 293},
		{"water", "residential", "Residential: 5-20 m3", 5, 20, "m3", 590},
		{"water", "residential", "Residential: 20-50 m3", 20, 50, "m3", 781},
		{"water", "residential", "Residential: >50 m3", 50, 0, "m3", 879},
		{"water", "commercial", "Commercial / Industrial", 0, 0, "m3", 879},
	}

	for _, t := range tariffs {
		if _, err := tx.Exec(
			"INSERT INTO tariffs (utility, category, description, min_usage, max_usage, unit, price_per_unit, currency) VALUES (?, ?, ?, ?, ?, ?, ?, 'RWF')",
			t.Utility, t.Category, t.Description, t.MinUsage, t.MaxUsage, t.Unit, t.PricePerUnit,
		); err != nil {
			return err
		}
	}
	return nil
}
