package database

import "database/sql"

// RunMigrations creates all tables if they don't exist.
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS provinces (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS districts (
			id TEXT PRIMARY KEY,
			province_id TEXT NOT NULL,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (province_id) REFERENCES provinces(id),
			UNIQUE(province_id, name)
		)`,
		`CREATE TABLE IF NOT EXISTS sectors (
			id TEXT PRIMARY KEY,
			district_id TEXT NOT NULL,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (district_id) REFERENCES districts(id),
			UNIQUE(district_id, name)
		)`,
		`CREATE TABLE IF NOT EXISTS cells (
			id TEXT PRIMARY KEY,
			sector_id TEXT NOT NULL,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (sector_id) REFERENCES sectors(id),
			UNIQUE(sector_id, name)
		)`,
		`CREATE TABLE IF NOT EXISTS villages (
			id TEXT PRIMARY KEY,
			cell_id TEXT NOT NULL,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (cell_id) REFERENCES cells(id),
			UNIQUE(cell_id, name)
		)`,
		`CREATE TABLE IF NOT EXISTS banks (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			short_name TEXT NOT NULL,
			swift_code TEXT NOT NULL UNIQUE,
			phone TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS branches (
			id TEXT PRIMARY KEY,
			bank_id TEXT NOT NULL,
			name TEXT NOT NULL,
			location TEXT NOT NULL DEFAULT '',
			phone TEXT NOT NULL DEFAULT '',
			FOREIGN KEY (bank_id) REFERENCES banks(id)
		)`,
		`CREATE TABLE IF NOT EXISTS tariffs (
			id TEXT PRIMARY KEY,
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
		`CREATE TABLE IF NOT EXISTS exchange_rates (
			id TEXT PRIMARY KEY,
			currency_code TEXT NOT NULL,
			currency_name TEXT NOT NULL,
			buying REAL NOT NULL,
			selling REAL NOT NULL,
			middle REAL NOT NULL,
			date TEXT NOT NULL,
			fetched_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS telecoms (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			short_name TEXT NOT NULL,
			type TEXT NOT NULL,
			website TEXT NOT NULL DEFAULT '',
			customer_care TEXT NOT NULL DEFAULT '',
			ussd_codes TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS telecom_plans (
			id TEXT PRIMARY KEY,
			telecom_id TEXT NOT NULL,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			speed TEXT NOT NULL DEFAULT '',
			data_cap TEXT NOT NULL DEFAULT '',
			price TEXT NOT NULL,
			validity TEXT NOT NULL DEFAULT '',
			FOREIGN KEY (telecom_id) REFERENCES telecoms(id)
		)`,
		`CREATE TABLE IF NOT EXISTS mobile_money (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			provider TEXT NOT NULL,
			ussd_code TEXT NOT NULL,
			agent_ussd TEXT NOT NULL DEFAULT '',
			daily_limit TEXT NOT NULL DEFAULT '',
			monthly_limit TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS mobile_money_fees (
			id TEXT PRIMARY KEY,
			service_id TEXT NOT NULL,
			tx_type TEXT NOT NULL,
			min_amount TEXT NOT NULL,
			max_amount TEXT NOT NULL,
			fee TEXT NOT NULL,
			FOREIGN KEY (service_id) REFERENCES mobile_money(id)
		)`,
		`CREATE TABLE IF NOT EXISTS emergency_numbers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			number TEXT NOT NULL,
			category TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS public_holidays (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			date TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			is_movable INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		// Indexes for common queries
		`CREATE INDEX IF NOT EXISTS idx_exchange_rates_currency ON exchange_rates(currency_code)`,
		`CREATE INDEX IF NOT EXISTS idx_districts_province ON districts(province_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sectors_district ON sectors(district_id)`,
		`CREATE INDEX IF NOT EXISTS idx_cells_sector ON cells(sector_id)`,
		`CREATE INDEX IF NOT EXISTS idx_villages_cell ON villages(cell_id)`,
		`CREATE INDEX IF NOT EXISTS idx_branches_bank ON branches(bank_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tariffs_utility ON tariffs(utility)`,
		`CREATE INDEX IF NOT EXISTS idx_telecom_plans_telecom ON telecom_plans(telecom_id)`,
		`CREATE INDEX IF NOT EXISTS idx_mobile_money_fees_service ON mobile_money_fees(service_id)`,
		`CREATE INDEX IF NOT EXISTS idx_emergency_numbers_category ON emergency_numbers(category)`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return err
		}
	}
	return nil
}
