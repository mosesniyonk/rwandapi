package handlers

import (
	"net/http"

	"github.com/mosesniyonk/rwandapi/internal/models"
)

func GetCountryInfo(w http.ResponseWriter, r *http.Request) {
	info := models.CountryInfo{
		Name:         "Rwanda",
		OfficialName: "Republic of Rwanda",
		Capital:      "Kigali",
		Population:   13246394,
		Area:         26338,
		Languages:    []string{"Kinyarwanda", "English", "French", "Swahili"},
		Currency:     "Rwandan Franc",
		CurrencyCode: "RWF",
		CallingCode:  "+250",
		TLD:          ".rw",
		Timezone:     "CAT (UTC+2)",
		DrivingSide:  "right",
		ISO2:         "RW",
		ISO3:         "RWA",
		Motto:        "Ubumwe, Umurimo, Gukunda Igihugu (Unity, Work, Patriotism)",
		Anthem:       "Rwanda Nziza (Beautiful Rwanda)",
		Government:   "Unitary presidential republic",
		President:    "Paul Kagame",
		Provinces:    5,
		Districts:    30,
		Sectors:      416,
	}

	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: info})
}
