package models

import "time"

// Province represents one of Rwanda's 5 provinces.
type Province struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Districts []District `json:"districts,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// District represents one of Rwanda's 30 districts.
type District struct {
	ID         string   `json:"id"`
	ProvinceID string   `json:"province_id"`
	Name       string   `json:"name"`
	Sectors    []Sector `json:"sectors,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Sector is a subdivision of a district.
type Sector struct {
	ID         string `json:"id"`
	DistrictID string `json:"district_id"`
	Name       string `json:"name"`
	Cells      []Cell `json:"cells,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Cell is a subdivision of a sector.
type Cell struct {
	ID       string    `json:"id"`
	SectorID string    `json:"sector_id"`
	Name     string    `json:"name"`
	Villages []Village `json:"villages,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Village is the smallest administrative unit.
type Village struct {
	ID        string    `json:"id"`
	CellID    string    `json:"cell_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Bank represents a bank operating in Rwanda.
type Bank struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ShortName string    `json:"short_name"`
	SwiftCode string    `json:"swift_code"`
	Phone     string    `json:"phone,omitempty"`
	Branches  []Branch  `json:"branches,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Branch represents a bank branch.
type Branch struct {
	ID       string `json:"id"`
	BankID   string `json:"bank_id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Phone    string `json:"phone,omitempty"`
}

// Tariff represents a utility tariff band.
type Tariff struct {
	ID          string  `json:"id"`
	Utility     string  `json:"utility"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	MinUsage    float64 `json:"min_usage"`
	MaxUsage    float64 `json:"max_usage"`
	Unit        string  `json:"unit"`
	PricePerUnit float64 `json:"price_per_unit"`
	Currency    string  `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
}

// ExchangeRate represents a BNR exchange rate for a currency against RWF.
type ExchangeRate struct {
	ID           string  `json:"id"`
	CurrencyCode string  `json:"currency_code"`
	CurrencyName string  `json:"currency_name"`
	Buying       float64 `json:"buying"`
	Selling      float64 `json:"selling"`
	Middle       float64 `json:"middle"`
	Date         string  `json:"date"`
	FetchedAt    string  `json:"fetched_at"`
}

// ExchangeRateResponse wraps exchange rate data with metadata.
type ExchangeRateResponse struct {
	Rates       []ExchangeRate `json:"rates"`
	LastUpdated string         `json:"last_updated"`
	Source      string         `json:"source"`
	Warning     string         `json:"warning,omitempty"`
}

// Telecom represents a telecom provider or ISP.
type Telecom struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	ShortName    string        `json:"short_name"`
	Type         string        `json:"type"`
	Website      string        `json:"website,omitempty"`
	CustomerCare string        `json:"customer_care,omitempty"`
	USSDCodes    []USSDCode    `json:"ussd_codes,omitempty"`
	Plans        []TelecomPlan `json:"plans,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
}

type USSDCode struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type TelecomPlan struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Speed    string `json:"speed,omitempty"`
	DataCap  string `json:"data_cap,omitempty"`
	Price    string `json:"price"`
	Validity string `json:"validity"`
}

// MobileMoneyService represents a mobile money provider.
type MobileMoneyService struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Provider     string           `json:"provider"`
	USSDCode     string           `json:"ussd_code"`
	AgentUSSD    string           `json:"agent_ussd,omitempty"`
	DailyLimit   string           `json:"daily_limit"`
	MonthlyLimit string           `json:"monthly_limit"`
	Fees         []MobileMoneyFee `json:"fees,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
}

type MobileMoneyFee struct {
	ID        string `json:"id"`
	TxType    string `json:"tx_type"`
	MinAmount string `json:"min_amount"`
	MaxAmount string `json:"max_amount"`
	Fee       string `json:"fee"`
}

// EmergencyNumber represents a useful phone number.
type EmergencyNumber struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Number      string `json:"number"`
	Category    string `json:"category"`
	Description string `json:"description,omitempty"`
}

// PublicHoliday represents a national holiday.
type PublicHoliday struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Date        string `json:"date"`
	Description string `json:"description,omitempty"`
	IsMovable   bool   `json:"is_movable"`
}

// CountryInfo contains metadata about Rwanda.
type CountryInfo struct {
	Name         string   `json:"name"`
	OfficialName string   `json:"official_name"`
	Capital      string   `json:"capital"`
	Population   int64    `json:"population"`
	Area         float64  `json:"area_sq_km"`
	Languages    []string `json:"languages"`
	Currency     string   `json:"currency"`
	CurrencyCode string   `json:"currency_code"`
	CallingCode  string   `json:"calling_code"`
	TLD          string   `json:"tld"`
	Timezone     string   `json:"timezone"`
	DrivingSide  string   `json:"driving_side"`
	ISO2         string   `json:"iso_alpha2"`
	ISO3         string   `json:"iso_alpha3"`
	Motto        string   `json:"motto"`
	Anthem       string   `json:"anthem"`
	Government   string   `json:"government"`
	President    string   `json:"president"`
	Provinces    int      `json:"provinces"`
	Districts    int      `json:"districts"`
	Sectors      int      `json:"sectors"`
}

// SearchResult represents a search hit across geographic entities.
type SearchResult struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

// APIResponse wraps all API responses.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Count   int         `json:"count,omitempty"`
}
