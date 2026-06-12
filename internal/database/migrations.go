package database

import "database/sql"

// RunMigrations creates all tables if they don't exist.
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS provinces (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS districts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			province_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (province_id) REFERENCES provinces(id),
			UNIQUE(province_id, name)
		)`,
		`CREATE TABLE IF NOT EXISTS sectors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			district_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (district_id) REFERENCES districts(id),
			UNIQUE(district_id, name)
		)`,
		`CREATE TABLE IF NOT EXISTS cells (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sector_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (sector_id) REFERENCES sectors(id),
			UNIQUE(sector_id, name)
		)`,
		`CREATE TABLE IF NOT EXISTS villages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			cell_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (cell_id) REFERENCES cells(id),
			UNIQUE(cell_id, name)
		)`,
		`CREATE TABLE IF NOT EXISTS banks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			short_name TEXT NOT NULL,
			swift_code TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS branches (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bank_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			location TEXT NOT NULL DEFAULT '',
			FOREIGN KEY (bank_id) REFERENCES banks(id)
		)`,
		`CREATE TABLE IF NOT EXISTS tariffs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			utility TEXT NOT NULL,
			category TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			min_usage REAL NOT NULL DEFAULT 0,
			max_usage REAL NOT NULL DEFAULT 0,
			unit TEXT NOT NULL,
			price_per_unit REAL NOT NULL,
			currency TEXT NOT NULL DEFAULT 'RWF',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		// Indexes for common queries
		`CREATE INDEX IF NOT EXISTS idx_districts_province ON districts(province_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sectors_district ON sectors(district_id)`,
		`CREATE INDEX IF NOT EXISTS idx_cells_sector ON cells(sector_id)`,
		`CREATE INDEX IF NOT EXISTS idx_villages_cell ON villages(cell_id)`,
		`CREATE INDEX IF NOT EXISTS idx_branches_bank ON branches(bank_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tariffs_utility ON tariffs(utility)`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return err
		}
	}
	return nil
}
