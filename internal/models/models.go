package models

import "time"

// Province represents one of Rwanda's 5 provinces.
type Province struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Districts []District `json:"districts,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// District represents one of Rwanda's 30 districts.
type District struct {
	ID         int      `json:"id"`
	ProvinceID int      `json:"province_id"`
	Name       string   `json:"name"`
	Sectors    []Sector `json:"sectors,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Sector is a subdivision of a district.
type Sector struct {
	ID         int    `json:"id"`
	DistrictID int    `json:"district_id"`
	Name       string `json:"name"`
	Cells      []Cell `json:"cells,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Cell is a subdivision of a sector.
type Cell struct {
	ID       int       `json:"id"`
	SectorID int       `json:"sector_id"`
	Name     string    `json:"name"`
	Villages []Village `json:"villages,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Village is the smallest administrative unit.
type Village struct {
	ID        int       `json:"id"`
	CellID    int       `json:"cell_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Bank represents a bank operating in Rwanda.
type Bank struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	ShortName string   `json:"short_name"`
	SwiftCode string   `json:"swift_code"`
	Branches  []Branch `json:"branches,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Branch represents a bank branch.
type Branch struct {
	ID       int    `json:"id"`
	BankID   int    `json:"bank_id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// Tariff represents a utility tariff band.
type Tariff struct {
	ID          int     `json:"id"`
	Utility     string  `json:"utility"`      // "electricity" or "water"
	Category    string  `json:"category"`      // "residential", "commercial", "industrial"
	Description string  `json:"description"`
	MinUsage    float64 `json:"min_usage"`     // kWh or m3
	MaxUsage    float64 `json:"max_usage"`     // kWh or m3, 0 means unlimited
	Unit        string  `json:"unit"`          // "kWh" or "m3"
	PricePerUnit float64 `json:"price_per_unit"` // RWF
	Currency    string  `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
}

// SearchResult represents a search hit across geographic entities.
type SearchResult struct {
	Type string      `json:"type"` // "province", "district", "sector", "cell", "village"
	ID   int         `json:"id"`
	Name string      `json:"name"`
	Path string      `json:"path"` // e.g. "Kigali City > Gasabo > Kimironko"
}

// APIResponse wraps all API responses.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Count   int         `json:"count,omitempty"`
}
