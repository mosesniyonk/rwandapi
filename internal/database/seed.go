package database

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
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
	if err := seedExchangeRates(tx); err != nil {
		return err
	}
	if err := seedTelecoms(tx); err != nil {
		return err
	}
	if err := seedMobileMoney(tx); err != nil {
		return err
	}
	if err := seedEmergencyNumbers(tx); err != nil {
		return err
	}
	if err := seedPublicHolidays(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func newID() string {
	return uuid.New().String()
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
		if _, err := tx.Exec("INSERT INTO provinces (id, name) VALUES (?, ?)", newID(), name); err != nil {
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
		var pid string
		if err := tx.QueryRow("SELECT id FROM provinces WHERE name = ?", province).Scan(&pid); err != nil {
			return err
		}
		for _, d := range dists {
			if _, err := tx.Exec("INSERT INTO districts (id, province_id, name) VALUES (?, ?, ?)", newID(), pid, d); err != nil {
				return err
			}
		}
	}
	return nil
}

func seedSectors(tx *sql.Tx) error {
	// Complete sectors for all 30 districts of Rwanda
	sectors := map[string][]string{
		// ── Kigali City ──
		"Gasabo": {
			"Bumbogo", "Gatsata", "Gikomero", "Gisozi", "Jabana",
			"Jali", "Kacyiru", "Kimihurura", "Kimironko", "Kinyinya",
			"Ndera", "Nduba", "Remera", "Rusororo", "Rutunga",
		},
		"Kicukiro": {
			"Gahanga", "Gatenga", "Gikondo", "Kagarama", "Kanombe",
			"Kicukiro", "Kigarama", "Masaka", "Niboye", "Nyarugunga",
		},
		"Nyarugenge": {
			"Gitega", "Kanyinya", "Kigali", "Kimisagara", "Mageragere",
			"Muhima", "Nyakabanda", "Nyamirambo", "Nyarugenge", "Rwezamenyo",
		},
		// ── Eastern Province ──
		"Bugesera": {
			"Gashora", "Juru", "Kamabuye", "Mareba", "Mayange",
			"Musenyi", "Mwogo", "Ngeruka", "Ntarama", "Nyamata",
			"Nyarugenge", "Rilima", "Ruhuha", "Rweru", "Shyara",
		},
		"Gatsibo": {
			"Gasange", "Gatsibo", "Gitoki", "Kabarore", "Kageyo",
			"Kiramuruzi", "Kiziguro", "Muhura", "Murambi", "Ngarama",
			"Nyagihanga", "Remera", "Rugarama", "Rwimbogo",
		},
		"Kayonza": {
			"Gahini", "Kabare", "Kabarondo", "Mukarange", "Murama",
			"Murundi", "Mwiri", "Ndego", "Nyamirama", "Rukara",
			"Ruramira", "Rwinkwavu",
		},
		"Kirehe": {
			"Gahara", "Gatore", "Kigarama", "Kigina", "Kirehe",
			"Mahama", "Mpanga", "Musaza", "Mushikiri", "Nasho",
			"Nyamugari", "Nyarubuye",
		},
		"Ngoma": {
			"Gashanda", "Jarama", "Karembo", "Kazo", "Kibungo",
			"Mugesera", "Murama", "Mutenderi", "Remera", "Rukira",
			"Rukumberi", "Rurenge", "Sake", "Zaza",
		},
		"Nyagatare": {
			"Gatunda", "Karama", "Karangazi", "Katabagemu", "Kiyombe",
			"Matimba", "Mimuri", "Mukama", "Musheri", "Nyagatare",
			"Rukomo", "Rwempasha", "Rwimiyaga", "Tabagwe",
		},
		"Rwamagana": {
			"Fumbwe", "Gahengeri", "Gishali", "Karenge", "Kigabiro",
			"Muhazi", "Munyaga", "Munyiginya", "Musha", "Muyumbu",
			"Mwulire", "Nyakaliro", "Nzige", "Rubona",
		},
		// ── Northern Province ──
		"Burera": {
			"Bungwe", "Butaro", "Cyanika", "Cyeru", "Gahunga",
			"Gatebe", "Gitovu", "Kagogo", "Kinoni", "Kinyababa",
			"Kivuye", "Nemba", "Rugarama", "Rugendabari", "Ruhunde",
			"Rusarabuye", "Rwerere",
		},
		"Gakenke": {
			"Busengo", "Coko", "Cyabingo", "Gakenke", "Gashenyi",
			"Janja", "Kamubuga", "Karambo", "Kivuruga", "Mataba",
			"Minazi", "Mugunga", "Muhondo", "Muyongwe", "Muzo",
			"Nemba", "Ruli", "Rusasa", "Rushashi",
		},
		"Gicumbi": {
			"Bukure", "Bwisige", "Byumba", "Cyumba", "Giti",
			"Kaniga", "Manyagiro", "Miyove", "Kageyo", "Mukarange",
			"Muko", "Mutete", "Nyamiyaga", "Nyankenke", "Rubaya",
			"Rukomo", "Rushaki", "Rutare", "Ruvune", "Rwamiko",
			"Shangasha",
		},
		"Musanze": {
			"Busogo", "Cyuve", "Gacaca", "Gashaki", "Gataraga",
			"Kimonyi", "Kinigi", "Muhoza", "Muko", "Musanze",
			"Nkotsi", "Nyange", "Remera", "Rwaza", "Shingiro",
		},
		"Rulindo": {
			"Base", "Burega", "Bushoki", "Buyoga", "Cyinzuzi",
			"Cyungo", "Kinihira", "Kisaro", "Masoro", "Mbogo",
			"Murambi", "Ngoma", "Ntarabana", "Rukozo", "Rusiga",
			"Shyorongi", "Tumba",
		},
		// ── Southern Province ──
		"Gisagara": {
			"Gikonko", "Gishubi", "Kansi", "Kibirizi", "Kigembe",
			"Mamba", "Muganza", "Mugombwa", "Mukingo", "Musha",
			"Ndora", "Nyanza", "Save",
		},
		"Huye": {
			"Gishamvu", "Huye", "Karama", "Kigoma", "Kinazi",
			"Maraba", "Mbazi", "Mukura", "Ngoma", "Ruhashya",
			"Rusatira", "Rwaniro", "Simbi", "Tumba",
		},
		"Kamonyi": {
			"Gacurabwenge", "Karama", "Kayenzi", "Kayumbu", "Mugina",
			"Musambira", "Ngamba", "Nyamiyaga", "Nyarubaka", "Rugarika",
			"Rukoma", "Runda",
		},
		"Muhanga": {
			"Cyeza", "Kabacuzi", "Kibangu", "Kiyumba", "Muhanga",
			"Mushishiro", "Nyabinoni", "Nyamabuye", "Nyarusange",
			"Rongi", "Rugendabari", "Shyogwe",
		},
		"Nyamagabe": {
			"Buruhukiro", "Cyanika", "Gasaka", "Gatare", "Kaduha",
			"Kamegeri", "Kibirizi", "Kibumbwe", "Kitabi", "Mbazi",
			"Mugano", "Musange", "Musebeya", "Mushubi", "Nkomane",
			"Tare", "Uwinkingi",
		},
		"Nyanza": {
			"Busasamana", "Busoro", "Cyabakamyi", "Kibirizi", "Kigoma",
			"Mukingo", "Muyira", "Ntyazo", "Nyagisozi", "Rwabicuma",
		},
		"Nyaruguru": {
			"Busanze", "Cyahinda", "Kibeho", "Kivu", "Mata",
			"Muganza", "Munini", "Ngera", "Ngoma", "Nyabimata",
			"Nyagisozi", "Ruheru", "Ruramba", "Rusenge",
		},
		"Ruhango": {
			"Bweramana", "Byimana", "Kabagali", "Kinazi", "Kinihira",
			"Mbuye", "Mwendo", "Ntongwe", "Ruhango",
		},
		// ── Western Province ──
		"Karongi": {
			"Bwishyura", "Gishari", "Gishyita", "Gitesi", "Mubuga",
			"Murambi", "Murundi", "Mutuntu", "Rubengera", "Rugabano",
			"Ruganda", "Rwankuba", "Twumba",
		},
		"Ngororero": {
			"Bwira", "Gatumba", "Hindiro", "Kabaya", "Kageyo",
			"Kavumu", "Matyazo", "Muhanda", "Muhororo", "Ndaro",
			"Ngororero", "Nyange", "Sovu",
		},
		"Nyabihu": {
			"Bigogwe", "Jenda", "Jomba", "Kabatwa", "Karago",
			"Kintobo", "Mukamira", "Muringa", "Rambura", "Rugera",
			"Rurembo", "Shyira",
		},
		"Nyamasheke": {
			"Bushekeri", "Bushenge", "Cyato", "Gihombo", "Kagano",
			"Kanjongo", "Karambi", "Karengera", "Kirimbi", "Macuba",
			"Mahembe", "Nyabitekeri", "Rangiro", "Ruharambuga", "Shangi",
		},
		"Rubavu": {
			"Bugeshi", "Busasamana", "Cyanzarwe", "Gisenyi", "Kanama",
			"Kanzenze", "Mudende", "Nyakiriba", "Nyamyumba", "Nyundo",
			"Rubavu", "Rugerero",
		},
		"Rusizi": {
			"Bugarama", "Butare", "Bweyeye", "Gashonga", "Giheke",
			"Gihundwe", "Gikundamvura", "Gitambi", "Kamembe", "Muganza",
			"Mururu", "Nkanka", "Nkombo", "Nkungu", "Nyakabuye",
			"Nyakarenzo", "Nzahaha", "Rwimbogo",
		},
		"Rutsiro": {
			"Boneza", "Gihango", "Kigeyo", "Kivumu", "Manihira",
			"Mukura", "Murunda", "Musasa", "Mushonyi", "Mushubati",
			"Nyabirasi", "Ruhango", "Rusebeya",
		},
	}

	for district, secs := range sectors {
		var did string
		if err := tx.QueryRow("SELECT id FROM districts WHERE name = ?", district).Scan(&did); err != nil {
			return err
		}
		for _, s := range secs {
			if _, err := tx.Exec("INSERT INTO sectors (id, district_id, name) VALUES (?, ?, ?)", newID(), did, s); err != nil {
				return err
			}
		}
	}
	return nil
}

func seedCells(tx *sql.Tx) error {
	type cellData struct {
		district string
		sector   string
		cells    []string
	}

	data := []cellData{
		// ════════════════════════════════════════
		// KIGALI CITY
		// ════════════════════════════════════════
		{"Gasabo", "Bumbogo", []string{"Kinyaga", "Musave", "Mvuzo", "Ngara", "Nkuzuzu", "Nyabikenke", "Nyagasozi"}},
		{"Gasabo", "Gatsata", []string{"Karuruma", "Nyamabuye", "Nyamugari"}},
		{"Gasabo", "Gikomero", []string{"Gasagara", "Gicaca", "Kibara", "Munini", "Murambi"}},
		{"Gasabo", "Gisozi", []string{"Musezero", "Ruhango"}},
		{"Gasabo", "Jabana", []string{"Akamatamu", "Bweramvura", "Kabuye", "Kidashya", "Ngiryi"}},
		{"Gasabo", "Jali", []string{"Agateko", "Buhiza", "Muko", "Nkusi", "Nyabuliba", "Nyakabungo", "Nyamitanga"}},
		{"Gasabo", "Kacyiru", []string{"Kamatamu", "Kamutwa", "Kibaza"}},
		{"Gasabo", "Kimihurura", []string{"Kamukina", "Kimihurura", "Rugando"}},
		{"Gasabo", "Kimironko", []string{"Bibare", "Kibagabaga", "Nyagatovu"}},
		{"Gasabo", "Kinyinya", []string{"Gacuriro", "Gasharu", "Kagugu", "Murama"}},
		{"Gasabo", "Ndera", []string{"Bwiza", "Cyaruzinge", "Kibenga", "Masoro", "Mukuyu", "Rudashya"}},
		{"Gasabo", "Nduba", []string{"Butare", "Gasanze", "Gasura", "Gatunga", "Muremure", "Sha", "Shango"}},
		{"Gasabo", "Remera", []string{"Nyabisindu", "Nyarutarama", "RukiriI", "RukiriII"}},
		{"Gasabo", "Rusororo", []string{"Bisenga", "Gasagara", "KabugaI", "KabugaII", "Kinyana", "Mbandazi", "Nyagahinga", "Ruhanga"}},
		{"Gasabo", "Rutunga", []string{"Gasabo", "Indatemwa", "Kabaliza", "Kacyatwa", "Kibenga", "Kigabiro"}},
		{"Kicukiro", "Gahanga", []string{"Gahanga", "Kagasa", "Karembure", "Murinja", "Nunga", "Rwabutenge"}},
		{"Kicukiro", "Gatenga", []string{"Gatenga", "Karambo", "Nyanza", "Nyarurama"}},
		{"Kicukiro", "Gikondo", []string{"Kagunga", "Kanserege", "Kinunga"}},
		{"Kicukiro", "Kagarama", []string{"Kanserege", "Muyange", "Rukatsa"}},
		{"Kicukiro", "Kanombe", []string{"Busanza", "Kabeza", "Karama", "Rubirizi"}},
		{"Kicukiro", "Kicukiro", []string{"Gasharu", "Kagina", "Kicukiro", "Ngoma"}},
		{"Kicukiro", "Kigarama", []string{"Bwerankori", "Karugira", "Kigarama", "Nyarurama", "Rwampara"}},
		{"Kicukiro", "Masaka", []string{"Ayabaraya", "Cyimo", "Gako", "Gasharu", "Gikomero", "Gishaka"}},
		{"Kicukiro", "Niboye", []string{"Gatare", "Niboye", "Nyakabanda"}},
		{"Kicukiro", "Nyarugunga", []string{"Kamashashi", "Nonko", "Rwimbogo"}},
		{"Nyarugenge", "Gitega", []string{"Akabahizi", "Akabeza", "Gacyamo", "Kigarama", "Kinyange", "Kora"}},
		{"Nyarugenge", "Kanyinya", []string{"Nyamweru", "Nzove", "Taba"}},
		{"Nyarugenge", "Kigali", []string{"Kigali", "Mwendo", "Nyabugogo", "Ruriba", "Rwesero"}},
		{"Nyarugenge", "Kimisagara", []string{"Kamuhoza", "Katabaro", "Kimisagara"}},
		{"Nyarugenge", "Mageragere", []string{"Kankuba", "Kavumu", "Mataba", "Ntungamo", "Nyarufunzo", "Nyarurenzi", "Runzenze"}},
		{"Nyarugenge", "Muhima", []string{"Amahoro", "Kabasengerezi", "Kabeza", "Nyabugogo", "Rugenge", "Tetero", "Ubumwe"}},
		{"Nyarugenge", "Nyakabanda", []string{"MunaniraI", "MunaniraII", "NyakabandaI", "NyakabandaII"}},
		{"Nyarugenge", "Nyamirambo", []string{"Cyivugiza", "Gasharu", "Mumena", "Rugarama"}},
		{"Nyarugenge", "Nyarugenge", []string{"Agatare", "Biryogo", "Kiyovu", "Rwampara"}},
		{"Nyarugenge", "Rwezamenyo", []string{"KabuguruI", "KabuguruII", "RwezamenyoI", "RwezamenyoII"}},

		// ════════════════════════════════════════
		// EASTERN PROVINCE
		// ════════════════════════════════════════
		{"Bugesera", "Gashora", []string{"Biryogo", "Kabuye", "Kagomasi", "Mwendo", "Ramiro"}},
		{"Bugesera", "Juru", []string{"Juru", "Kabukuba", "Mugorore", "Musovu", "Rwinume"}},
		{"Bugesera", "Kamabuye", []string{"Biharagu", "Burenge", "Kampeka", "Nyakayaga", "Tunda"}},
		{"Bugesera", "Mareba", []string{"Bushenyi", "Gakomeye", "Nyamigina", "Rango", "Rugarama"}},
		{"Bugesera", "Mayange", []string{"Gakamba", "Kagenge", "Kibenga", "Kibirizi", "Mbyo"}},
		{"Bugesera", "Musenyi", []string{"Gicaca", "Musenyi", "Nyagihunika", "Rulindo"}},
		{"Bugesera", "Mwogo", []string{"Bitaba", "Kagasa", "Rugunga", "Rurenge"}},
		{"Bugesera", "Ngeruka", []string{"Gihembe", "Murama", "Ngeruka", "Nyakayenzi", "Rutonde"}},
		{"Bugesera", "Ntarama", []string{"Cyugaro", "Kanzenze", "Kibungo"}},
		{"Bugesera", "Nyamata", []string{"Kanazi", "Kayumba", "Maranyundo", "Murama", "NyamataYUmujyi"}},
		{"Bugesera", "Nyarugenge", []string{"Gihinga", "Kabuye", "Murambi", "Ngenda", "Rugando"}},
		{"Bugesera", "Rilima", []string{"Kabeza", "Karera", "Kimaranzara", "Ntarama", "Nyabagendwa"}},
		{"Bugesera", "Ruhuha", []string{"Bihari", "Gatanga", "Gikundamvura", "Kindama", "Ruhuha"}},
		{"Bugesera", "Rweru", []string{"Batima", "Kintambwe", "Mazane", "Nemba", "Nkanga", "Sharita"}},
		{"Bugesera", "Shyara", []string{"Kabagugu", "Kamabuye", "Nziranziza", "Rebero", "Rutare"}},
		{"Gatsibo", "Gasange", []string{"Kigabiro", "Kimana", "Teme", "Viro"}},
		{"Gatsibo", "Gatsibo", []string{"Gatsibo", "Manishya", "Mugera", "Nyabicwamba", "Nyagahanga"}},
		{"Gatsibo", "Gitoki", []string{"Bukomane", "Cyabusheshe", "Karubungo", "Mpondwa", "Nyamirama", "Rubira"}},
		{"Gatsibo", "Kabarore", []string{"Kabarore", "Kabeza", "Karenge", "Marimba", "Nyabikiri", "Simbwa"}},
		{"Gatsibo", "Kageyo", []string{"Busetsa", "Gituza", "Kintu", "Nyagisozi"}},
		{"Gatsibo", "Kiramuruzi", []string{"Akabuga", "Gakenke", "Gakoni", "Nyabisindu"}},
		{"Gatsibo", "Kiziguro", []string{"Agakomeye", "Mbogo", "Ndatemwa", "Rubona"}},
		{"Gatsibo", "Muhura", []string{"Bibare", "Gakorokombe", "Mamfu", "Rumuli", "Taba"}},
		{"Gatsibo", "Murambi", []string{"Murambi", "Nyamiyaga", "Rwankuba", "Rwimitereri"}},
		{"Gatsibo", "Ngarama", []string{"Bugamba", "Karambi", "Kigasha", "Ngarama", "Nyarubungo"}},
		{"Gatsibo", "Nyagihanga", []string{"Gitinda", "Kibare", "Mayange", "Murambi", "Nyagitabire", "Nyamirama"}},
		{"Gatsibo", "Remera", []string{"Bushobora", "Butiruka", "Kigabiro", "Nyagakombe", "Rurenge", "Rwarenga"}},
		{"Gatsibo", "Rugarama", []string{"Bugarama", "Gihuta", "Kanyangese", "Matare", "Matunguru", "Remera"}},
		{"Gatsibo", "Rwimbogo", []string{"Kiburara", "Munini", "Nyamatete", "Rwikiniro"}},
		{"Kayonza", "Gahini", []string{"Juru", "Kahi", "Kiyenzi", "Urugarama"}},
		{"Kayonza", "Kabare", []string{"Cyarubare", "Gitara", "Kirehe", "Rubimba", "Rubumba"}},
		{"Kayonza", "Kabarondo", []string{"Cyabajwa", "Cyinzovu", "Kabura", "Rusera"}},
		{"Kayonza", "Mukarange", []string{"Bwiza", "Kayonza", "Mburabuturo", "Nyagatovu", "Rugendabari"}},
		{"Kayonza", "Murama", []string{"Bunyentongo", "Muko", "Murama", "Nyakanazi", "Rusave"}},
		{"Kayonza", "Murundi", []string{"Buhabwa", "Karambi", "Murundi", "Ryamanyoni"}},
		{"Kayonza", "Mwiri", []string{"Kageyo", "Migera", "Nyamugari", "Nyawera"}},
		{"Kayonza", "Ndego", []string{"Byimana", "Isangano", "Karambi", "Kiyovu"}},
		{"Kayonza", "Nyamirama", []string{"Gikaya", "Musumba", "Rurambi", "Shyogo"}},
		{"Kayonza", "Rukara", []string{"Kawangire", "Rukara", "Rwimishinya"}},
		{"Kayonza", "Ruramira", []string{"Bugambira", "Nkamba", "Ruyonza", "Umubuga"}},
		{"Kayonza", "Rwinkwavu", []string{"Gihinga", "Mbarara", "Mukoyoyo", "Nkondo"}},
		{"Kirehe", "Gahara", []string{"Butezi", "Muhamba", "Murehe", "Nyagasenyi", "Nyakagezi", "Rubimba"}},
		{"Kirehe", "Gatore", []string{"Curazo", "Cyunuzi", "Muganza", "Nyamiryango", "Rwabutazi", "Rwantonde"}},
		{"Kirehe", "Kigarama", []string{"Cyanya", "Kigarama", "Kiremera", "Nyakerera", "Nyankurazo"}},
		{"Kirehe", "Kigina", []string{"Gatarama", "Rugarama", "Ruhanga", "Rwanteru"}},
		{"Kirehe", "Kirehe", []string{"Gahama", "Kirehe", "Nyabigega", "Nyabikokora", "Rwesero"}},
		{"Kirehe", "Mahama", []string{"Kamombo", "Munini", "Mwoga", "Saruhembe", "Umunini"}},
		{"Kirehe", "Mpanga", []string{"Bwiyorere", "Kankobwa", "Mpanga", "Mushongi", "Nasho", "Nyakabungo", "Rubaya"}},
		{"Kirehe", "Musaza", []string{"Gasarabwayi", "Kabuga", "Mubuga", "Musaza", "Nganda"}},
		{"Kirehe", "Mushikiri", []string{"Bisagara", "Cyamigurwa", "Rugarama", "Rwanyamuhanga", "Rwayikona"}},
		{"Kirehe", "Nasho", []string{"Cyambwe", "Kagese", "Ntaruka", "Rubirizi", "Rugoma"}},
		{"Kirehe", "Nyamugari", []string{"Bukora", "Kagasa", "Kazizi", "Kiyanzi", "Nyamugari"}},
		{"Kirehe", "Nyarubuye", []string{"Mareba", "Nyabitare", "Nyarutunga"}},
		{"Ngoma", "Gashanda", []string{"Cyerwa", "Giseri", "Munege", "Mutsindo"}},
		{"Ngoma", "Jarama", []string{"Ihanika", "Jarama", "Karenge", "Kibimba", "Kigoma"}},
		{"Ngoma", "Karembo", []string{"Akaziba", "Karaba", "Nyamirambo"}},
		{"Ngoma", "Kazo", []string{"Birenga", "Gahurire", "Karama", "Kinyonzo", "Umukamba"}},
		{"Ngoma", "Kibungo", []string{"Cyasemakamba", "Gahima", "Gatonde", "Karenge", "Mahango"}},
		{"Ngoma", "Mugesera", []string{"Akabungo", "Mugatare", "Ntanga", "Nyamugari", "Nyange"}},
		{"Ngoma", "Murama", []string{"Gitaraga", "Kigabiro", "Mvumba", "Rurenge", "Sakara"}},
		{"Ngoma", "Mutenderi", []string{"Karwema", "Kibare", "Mutenderi", "Muzingira", "Nyagasozi"}},
		{"Ngoma", "Remera", []string{"Bugera", "Kinunga", "Ndekwe", "Nyamagana"}},
		{"Ngoma", "Rukira", []string{"Buliba", "Kibatsi", "Nyaruvumu", "Nyinya"}},
		{"Ngoma", "Rukumberi", []string{"Gituza", "Ntovi", "Rubago", "Rubona", "Rwintashya"}},
		{"Ngoma", "Rurenge", []string{"Akagarama", "Muhurire", "Musya", "Rugese", "Rujambara", "Rwikubo"}},
		{"Ngoma", "Sake", []string{"Gafunzo", "Kibonde", "Nkanga", "Rukoma"}},
		{"Ngoma", "Zaza", []string{"Nyagasozi", "Nyagatugunda", "Ruhembe", "Ruhinga"}},
		{"Nyagatare", "Gatunda", []string{"Cyagaju", "Kabeza", "Nyamikamba", "Nyamirembe", "Nyangara", "Nyarurema", "Rwensheke"}},
		{"Nyagatare", "Karama", []string{"Bushara", "Cyenkwanzi", "Gikagati", "Gikundamvura", "Kabuga", "Ndego", "Nyakiga"}},
		{"Nyagatare", "Karangazi", []string{"Kamate", "Karama", "Kizirakome", "Mbare", "Musenyi", "Ndama", "Nyagashanga", "Nyamirama", "Rubagabaga", "Rwenyemera", "Rwisirabo"}},
		{"Nyagatare", "Katabagemu", []string{"Bayigaburire", "Kaduha", "Kanyeganyege", "Katabagemu", "Kigarama", "Nyakigando", "Rubira", "Rugazi", "Rutoma"}},
		{"Nyagatare", "Kiyombe", []string{"Gataba", "Gitenga", "Kabungo", "Karambo", "Karujumba", "Tovu"}},
		{"Nyagatare", "Matimba", []string{"Bwera", "Byimana", "Cyembogo", "Kagitumba", "Kanyonza", "Matimba", "Nyabwishongwezi", "Rwentanga"}},
		{"Nyagatare", "Mimuri", []string{"Bibare", "Gakoma", "Mahoro", "Mimuri", "Rugari"}},
		{"Nyagatare", "Mukama", []string{"Bufunda", "Gatete", "Gihengeri", "Gishororo", "Kagina", "Rugarama"}},
		{"Nyagatare", "Musheri", []string{"Kibirizi", "Kijojo", "Musheri", "Ntoma", "Nyagatabire", "Nyamiyonga", "RugaramaI", "RugaramaII"}},
		{"Nyagatare", "Nyagatare", []string{"Barija", "Bushoga", "Cyabayaga", "Gakirage", "Kamagiri", "Nsheke", "Nyagatare", "Rutaraka", "Ryabega"}},
		{"Nyagatare", "Rukomo", []string{"Gahurura", "Gashenyi", "Nyakagarama", "RukomoII", "Rurenge"}},
		{"Nyagatare", "Rwempasha", []string{"Cyenjonjo", "Gasinga", "Kabare", "Kazaza", "Mishenyi", "Rugarama", "Rukorota", "Rutare", "Rwempasha", "Ryeru"}},
		{"Nyagatare", "Rwimiyaga", []string{"Gacundezi", "Kabeza", "Kirebe", "Ntoma", "Nyarupfubire", "Nyendo", "Rutungu", "Rwimiyaga"}},
		{"Nyagatare", "Tabagwe", []string{"Gishuro", "Gitengure", "Nkoma", "Nyabitekeri", "Nyagatoma", "Shonga", "Tabagwe"}},
		{"Rwamagana", "Fumbwe", []string{"Mununu", "Nyagasambu", "Nyakagunga", "Nyamirama", "Nyarubuye", "Sasabirago"}},
		{"Rwamagana", "Gahengeri", []string{"Gihumuza", "Kagezi", "Kanyangese", "Kibare", "Mutamwa", "Rugarama", "Runyinya", "Rweri"}},
		{"Rwamagana", "Gishali", []string{"Binunga", "Bwinsanga", "Cyinyana", "Gati", "Kavumu", "Ruhimbi", "Ruhunda"}},
		{"Rwamagana", "Karenge", []string{"Bicaca", "Byimana", "Kabasore", "Kangamba", "Karenge", "Nyabubare", "Nyamatete"}},
		{"Rwamagana", "Kigabiro", []string{"Bwiza", "Cyanya", "Nyagasenyi", "Sibagire", "Sovu"}},
		{"Rwamagana", "Muhazi", []string{"Byeza", "Kabare", "Karambi", "Karitutu", "Kitazigurwa", "Murambi", "Nsinda", "Ntebe", "Nyarusange"}},
		{"Rwamagana", "Munyaga", []string{"Kaduha", "Nkungu", "Rweru", "Zinga"}},
		{"Rwamagana", "Munyiginya", []string{"Binunga", "Bwana", "Cyarukamba", "Cyimbazi", "Nkomangwa", "Nyarubuye"}},
		{"Rwamagana", "Musha", []string{"Akabare", "Budahanda", "Kagarama", "Musha", "Nyabisindu", "Nyakabanda"}},
		{"Rwamagana", "Muyumbu", []string{"Akinyambo", "Bujyujyu", "Murehe", "Ntebe", "Nyarukombe"}},
		{"Rwamagana", "Mwulire", []string{"Bicumbi", "Bushenyi", "Mwulire", "Ntunga"}},
		{"Rwamagana", "Nyakaliro", []string{"Bihembe", "Gatare", "Gishore", "Munini", "Rwimbogo"}},
		{"Rwamagana", "Nzige", []string{"Akanzu", "Kigarama", "Murama", "Rugarama"}},
		{"Rwamagana", "Rubona", []string{"Byinza", "Kabatasi", "Kabuye", "Karambi", "Mabare", "Nawe"}},

		// ════════════════════════════════════════
		// NORTHERN PROVINCE
		// ════════════════════════════════════════
		{"Burera", "Bungwe", []string{"Bungwe", "Bushenya", "Mudugari", "Tumba"}},
		{"Burera", "Butaro", []string{"Gatsibo", "Mubuga", "Muhotora", "Nyamicucu", "Rusumo"}},
		{"Burera", "Cyanika", []string{"Gasiza", "Gisovu", "Kabyiniro", "Kagitega", "Kamanyana", "Nyagahinga"}},
		{"Burera", "Cyeru", []string{"Butare", "Ndongozi", "Ruyange"}},
		{"Burera", "Gahunga", []string{"Buramba", "Gisizi", "Kidakama", "Nyangwe", "Rwasa"}},
		{"Burera", "Gatebe", []string{"Gabiro", "Musenda", "Rwambogo", "Rwasa"}},
		{"Burera", "Gitovu", []string{"Mariba", "Musasa", "Runoga"}},
		{"Burera", "Kagogo", []string{"Kabaya", "Kayenzi", "Kiringa", "Nyamabuye"}},
		{"Burera", "Kinoni", []string{"Gafuka", "Nkenke", "Nkumba", "Ntaruka"}},
		{"Burera", "Kinyababa", []string{"Bugamba", "Kaganda", "Musasa", "Rutovu"}},
		{"Burera", "Kivuye", []string{"Bukwashuri", "Gashanje", "Murwa", "Nyirataba"}},
		{"Burera", "Nemba", []string{"Kivumu", "Nyamugari", "Rubona", "Rushara"}},
		{"Burera", "Rugarama", []string{"Cyahi", "Gafumba", "Karangara", "Rurembo"}},
		{"Burera", "Rugendabari", []string{"Kilibata", "Mucaca", "Nyanamo", "Rukandabyuma"}},
		{"Burera", "Ruhunde", []string{"Gaseke", "Gatare", "Gitovu", "Rusekera"}},
		{"Burera", "Rusarabuye", []string{"Kabona", "Ndago", "Ruhanga"}},
		{"Burera", "Rwerere", []string{"Gacundura", "Gashoro", "Ruconsho", "Rugari"}},
		{"Gakenke", "Busengo", []string{"Birambo", "Butereri", "Byibuhiro", "Kamina", "Kirabo", "Mwumba", "Ruhanga"}},
		{"Gakenke", "Coko", []string{"Kiruku", "Mbirima", "Nyange", "Nyanza"}},
		{"Gakenke", "Cyabingo", []string{"Muhaza", "Muhororo", "Muramba", "Mutanda", "Rukore"}},
		{"Gakenke", "Gakenke", []string{"Buheta", "Kagoma", "Nganzo", "Rusagara"}},
		{"Gakenke", "Gashenyi", []string{"Nyacyina", "Rukura", "Rutabo", "Rutenderi", "Taba"}},
		{"Gakenke", "Janja", []string{"Gakindo", "Gashyamba", "Gatwa", "Karukungu"}},
		{"Gakenke", "Kamubuga", []string{"Kamubuga", "Kidomo", "Mbatabata", "Rukore"}},
		{"Gakenke", "Karambo", []string{"Kanyanza", "Karambo", "Kirebe"}},
		{"Gakenke", "Kivuruga", []string{"Cyintare", "Gasiza", "Rugimbu", "Ruhinga", "Sereri"}},
		{"Gakenke", "Mataba", []string{"Buyange", "Gikombe", "Nyundo"}},
		{"Gakenke", "Minazi", []string{"Gasiho", "Munyana", "Murambi", "Raba"}},
		{"Gakenke", "Mugunga", []string{"Gahinga", "Munyana", "Mutego", "Nkomane", "Rutabo", "Rutenderi", "Rwamambe"}},
		{"Gakenke", "Muhondo", []string{"Busake", "Bwenda", "Gasiza", "Gihinga", "Huro", "Musagara", "Musenyi", "Ruganda", "Rwinkuba"}},
		{"Gakenke", "Muyongwe", []string{"Bumba", "Gisiza", "Karyango", "Nganzo", "Va"}},
		{"Gakenke", "Muzo", []string{"Kabatezi", "Kiryamo", "Mubuga", "Mwiyando", "Rwa"}},
		{"Gakenke", "Nemba", []string{"Buranga", "Gahinga", "Gisozi", "Mucaca"}},
		{"Gakenke", "Ruli", []string{"Busoro", "Gikingo", "Jango", "Ruli", "Rwesero"}},
		{"Gakenke", "Rusasa", []string{"Gataba", "Kamonyi", "Murambi", "Nyundo", "Rumbi", "Rurembo"}},
		{"Gakenke", "Rushashi", []string{"Burimba", "Busanane", "Joma", "Kageyo", "Mbogo", "Razi", "Rwankuba", "Shyombwe"}},
		{"Gicumbi", "Bukure", []string{"Karenge", "Kigabiro", "Kivumu", "Rwesero"}},
		{"Gicumbi", "Bwisige", []string{"Bwisige", "Gihuke", "Mukono", "Nyabushingitwa"}},
		{"Gicumbi", "Byumba", []string{"Gacurabwenge", "Gisuna", "Kibali", "Kivugiza", "Murama", "Ngondore", "Nyakabungo", "Nyamabuye", "Nyarutarama"}},
		{"Gicumbi", "Cyumba", []string{"Gasunzu", "Muhambo", "Nyakabungo", "Nyambare", "Nyaruka", "Rwankonjo"}},
		{"Gicumbi", "Giti", []string{"Gatobotobo", "Murehe", "Tanda"}},
		{"Gicumbi", "Kageyo", []string{"Gihembe", "Horezo", "Kabuga", "Muhondo", "Nyamiyaga"}},
		{"Gicumbi", "Kaniga", []string{"Bugomba", "Gatoma", "Mulindi", "Nyarwambu", "Rukurura"}},
		{"Gicumbi", "Manyagiro", []string{"Kabuga", "Nyiragifumba", "Nyiravugiza", "Remera", "Rusekera", "Ryaruyumba"}},
		{"Gicumbi", "Miyove", []string{"Gakenke", "Miyove", "Mubuga"}},
		{"Gicumbi", "Mukarange", []string{"Cyamuganga", "Gatenga", "Kiruhura", "Mutarama", "Rugerero", "Rusambya"}},
		{"Gicumbi", "Muko", []string{"Cyamuhinda", "Kigoma", "Mwendo", "Ngange", "Rebero"}},
		{"Gicumbi", "Mutete", []string{"Gaseke", "Kabeza", "Musenyi", "Mutandi", "Nyarubuye"}},
		{"Gicumbi", "Nyamiyaga", []string{"Gahumuliza", "Jamba", "Kabeza", "Kabuga", "Karambo", "Kiziba", "Mataba"}},
		{"Gicumbi", "Nyankenke", []string{"Butare", "Kigogo", "Kinishya", "Rusasa", "Rutete", "Rwagihura", "Yaramba"}},
		{"Gicumbi", "Rubaya", []string{"Gihanga", "Gishambashayo", "Gishari", "Muguramo", "Nyamiyaga"}},
		{"Gicumbi", "Rukomo", []string{"Cyeya", "Cyuru", "Gisiza", "Kinyami", "Mabare", "Munyinya"}},
		{"Gicumbi", "Rushaki", []string{"Gitega", "Kamutora", "Karurama"}},
		{"Gicumbi", "Rutare", []string{"Bikumba", "Gasharu", "Gatwaro", "Kigabiro", "Munanira", "Nkoto"}},
		{"Gicumbi", "Ruvune", []string{"Cyandaro", "Gasambya", "Gashirira", "Kabare", "Rebero", "Ruhondo"}},
		{"Gicumbi", "Rwamiko", []string{"Cyeru", "Kigabiro", "Nyagahinga"}},
		{"Gicumbi", "Shangasha", []string{"Bushara", "Kitazigurwa", "Nyabishambi", "Nyabubare", "Shangasha"}},
		{"Musanze", "Busogo", []string{"Gisesero", "Kavumu", "Nyagisozi", "Sahara"}},
		{"Musanze", "Cyuve", []string{"Bukinanyana", "Buruba", "Cyanya", "Kabeza", "Migeshi", "Rwebeya"}},
		{"Musanze", "Gacaca", []string{"Gakoro", "Gasakuza", "Kabirizi", "Karwasa"}},
		{"Musanze", "Gashaki", []string{"Kigabiro", "Kivumu", "Mbwe", "Muharuro"}},
		{"Musanze", "Gataraga", []string{"Mudakama", "Murago", "Rubindi", "Rungu"}},
		{"Musanze", "Kimonyi", []string{"Birira", "Buramira", "Kivumu", "Mbizi"}},
		{"Musanze", "Kinigi", []string{"Bisoke", "Kaguhu", "Kampanga", "Nyabigoma", "Nyonirima"}},
		{"Musanze", "Muhoza", []string{"Cyabararika", "Kigombe", "Mpenge", "Ruhengeri"}},
		{"Musanze", "Muko", []string{"Cyivugiza", "Cyogo", "Mburabuturo", "Songa"}},
		{"Musanze", "Musanze", []string{"Cyabagarura", "Garuka", "Kabazungu", "Nyarubuye", "Rwambogo"}},
		{"Musanze", "Nkotsi", []string{"Bikara", "Gashinga", "Mubago", "Rugeshi", "Ruyumba"}},
		{"Musanze", "Nyange", []string{"Cyivugiza", "Kabeza", "Kamwumba", "Muhabura", "Ninda"}},
		{"Musanze", "Remera", []string{"Gasongero", "Kamisave", "Murandi", "Murwa", "Rurambo"}},
		{"Musanze", "Rwaza", []string{"Bumara", "Kabushinge", "Musezero", "Nturo", "Nyarubuye"}},
		{"Musanze", "Shingiro", []string{"Gakingo", "Kibuguzo", "Mudende", "Mugari"}},
		{"Rulindo", "Base", []string{"Cyohoha", "Gitare", "Rwamahwa"}},
		{"Rulindo", "Burega", []string{"Butangampundu", "Karengeri", "Taba"}},
		{"Rulindo", "Bushoki", []string{"Gasiza", "Giko", "Kayenzi", "Mukoto", "Nyirangarama"}},
		{"Rulindo", "Buyoga", []string{"Busoro", "Butare", "Gahororo", "Gitumba", "Karama", "Mwumba", "Ndarage"}},
		{"Rulindo", "Cyinzuzi", []string{"Budakiranya", "Migendezo", "Rudogo"}},
		{"Rulindo", "Cyungo", []string{"Burehe", "Marembo", "Rwili"}},
		{"Rulindo", "Kinihira", []string{"Butunzi", "Karegamazi", "Marembo", "Rebero"}},
		{"Rulindo", "Kisaro", []string{"Gitatsa", "Kamushenyi", "Kigarama", "Mubuga", "Murama", "Sayo"}},
		{"Rulindo", "Masoro", []string{"Kabuga", "Kigarama", "Kivugiza", "Nyamyumba", "Shengampuli"}},
		{"Rulindo", "Mbogo", []string{"Bukoro", "Mushari", "Ngiramazi", "Rurenge"}},
		{"Rulindo", "Murambi", []string{"Bubangu", "Gatwa", "Mugambazi", "Mvuzo"}},
		{"Rulindo", "Ngoma", []string{"Kabuga", "Karambo", "Mugote", "Munyarwanda"}},
		{"Rulindo", "Ntarabana", []string{"Kajevuba", "Kiyanza", "Mahaza"}},
		{"Rulindo", "Rukozo", []string{"Buraro", "Bwimo", "Mberuka", "Mbuye"}},
		{"Rulindo", "Rusiga", []string{"Gako", "Kirenge", "Taba"}},
		{"Rulindo", "Shyorongi", []string{"Bugaragara", "Kijabagwe", "Muvumu", "Rubona", "Rutonde"}},
		{"Rulindo", "Tumba", []string{"Barari", "Gahabwa", "Misezero", "Nyirabirori", "Taba"}},

		// ════════════════════════════════════════
		// SOUTHERN PROVINCE
		// ════════════════════════════════════════
		{"Gisagara", "Gikonko", []string{"Cyiri", "Gasagara", "Gikonko", "Mbogo"}},
		{"Gisagara", "Gishubi", []string{"Gabiro", "Nyabitare", "Nyakibungo", "Nyeranzi"}},
		{"Gisagara", "Kansi", []string{"Akaboti", "Bwiza", "Sabusaro", "Umunini"}},
		{"Gisagara", "Kibirizi", []string{"Duwani", "Kibirizi", "Muyira", "Ruturo"}},
		{"Gisagara", "Kigembe", []string{"Agahabwa", "Gatovu", "Impinga", "Nyabikenke", "Rubona", "Rusagara"}},
		{"Gisagara", "Mamba", []string{"Gakoma", "Kabumbwe", "Mamba", "Muyaga", "Ramba"}},
		{"Gisagara", "Muganza", []string{"Cyumba", "Muganza", "Remera", "Rwamiko", "Saga"}},
		{"Gisagara", "Mugombwa", []string{"Baziro", "Kibayi", "Kibu", "Mugombwa", "Mukomacara"}},
		{"Gisagara", "Mukingo", []string{"Gitega", "Mukiza", "Nyabisagara", "Runyinya"}},
		{"Gisagara", "Musha", []string{"Bukinanyana", "Gatovu", "Kigarama", "Kimana"}},
		{"Gisagara", "Ndora", []string{"Bweya", "Cyamukuza", "Dahwe", "Gisagara", "Mukande"}},
		{"Gisagara", "Nyanza", []string{"Higiro", "Nyamugari", "Nyaruteja", "Umubanga"}},
		{"Gisagara", "Save", []string{"Gatoki", "Munazi", "Rwanza", "Shyanda", "Zivu"}},
		{"Huye", "Gishamvu", []string{"Nyakibanda", "Nyumba", "Ryakibogo", "Shori"}},
		{"Huye", "Huye", []string{"Muyogoro", "Nyakagezi", "Rukira", "Sovu"}},
		{"Huye", "Karama", []string{"Buhoro", "Bunazi", "Gahororo", "Kibingo", "Muhembe"}},
		{"Huye", "Kigoma", []string{"Gishihe", "Kabatwa", "Kabuga", "Karambi", "Musebeya", "Nyabisindu", "Rugarama", "Shanga"}},
		{"Huye", "Kinazi", []string{"Byinza", "Gahana", "Gitovu", "Kabona", "Sazange"}},
		{"Huye", "Maraba", []string{"Buremera", "Gasumba", "Kabuye", "Kanyinya", "Shanga", "Shyembe"}},
		{"Huye", "Mbazi", []string{"Gatobotobo", "Kabuga", "Mutunda", "Mwulire", "Rugango", "Rusagara", "Tare"}},
		{"Huye", "Mukura", []string{"Bukomeye", "Buvumu", "Icyeru", "RangoA"}},
		{"Huye", "Ngoma", []string{"Butare", "Kaburemera", "Matyazo", "Ngoma"}},
		{"Huye", "Ruhashya", []string{"Busheshi", "Gatovu", "Karama", "Mara", "Muhororo", "Rugogwe", "Ruhashya"}},
		{"Huye", "Rusatira", []string{"Buhimba", "Gafumba", "Kimirehe", "Kimuna", "Kiruhura", "Mugogwe"}},
		{"Huye", "Rwaniro", []string{"Gatwaro", "Kamwambi", "Kibiraro", "Mwendo", "Nyamabuye", "Nyaruhombo", "Shyunga"}},
		{"Huye", "Simbi", []string{"Cyendajuru", "Gisakura", "Kabusanza", "Mugobore", "Nyangazi"}},
		{"Huye", "Tumba", []string{"Cyarwa", "Cyimana", "Gitwa", "Mpare", "RangoB"}},
		{"Kamonyi", "Gacurabwenge", []string{"Gihinga", "Gihira", "Kigembe", "Nkingo"}},
		{"Kamonyi", "Karama", []string{"Bitare", "Bunyonga", "Muganza", "Nyamirembe"}},
		{"Kamonyi", "Kayenzi", []string{"Bugarama", "Cubi", "Kayonza", "Kirwa", "Mataba", "Nyamirama"}},
		{"Kamonyi", "Kayumbu", []string{"Busoro", "Gaseke", "Giko", "Muyange"}},
		{"Kamonyi", "Mugina", []string{"Jenda", "Kabugondo", "Mbati", "Mugina", "Nteko"}},
		{"Kamonyi", "Musambira", []string{"Buhoro", "Cyambwe", "Karengera", "Kivumu", "Mpushi", "Rukambura"}},
		{"Kamonyi", "Ngamba", []string{"Kabuga", "Kazirabonde", "Marembo"}},
		{"Kamonyi", "Nyamiyaga", []string{"Bibungo", "Kabashumba", "Kidahwe", "Mukinga", "Ngoma"}},
		{"Kamonyi", "Nyarubaka", []string{"Gitare", "Kambyeyi", "Kigusa", "Nyagishubi", "Ruyanza"}},
		{"Kamonyi", "Rugarika", []string{"Bihembe", "Kigese", "Masaka", "Nyarubuye", "Sheli"}},
		{"Kamonyi", "Rukoma", []string{"Bugoba", "Buguri", "Gishyeshye", "Murehe", "Mwirute", "Remera", "Taba"}},
		{"Kamonyi", "Runda", []string{"Gihara", "Kabagesera", "Kagina", "Muganza", "Ruyenzi"}},
		{"Muhanga", "Cyeza", []string{"Biringaga", "Kigarama", "Kivumu", "Makera", "Nyarunyinya", "Shori"}},
		{"Muhanga", "Kabacuzi", []string{"Buramba", "Butare", "Kabuye", "Kavumu", "Kibyimba", "Ngarama", "Ngoma", "Sholi"}},
		{"Muhanga", "Kibangu", []string{"Gisharu", "Gitega", "Jurwe", "Mubuga", "Rubyiniro", "Ryakanimba"}},
		{"Muhanga", "Kiyumba", []string{"Budende", "Ndago", "Remera", "Ruhina", "Rukeri"}},
		{"Muhanga", "Muhanga", []string{"Kanyinya", "Nganzo", "Nyamirama", "Remera", "Tyazo"}},
		{"Muhanga", "Mushishiro", []string{"Matyazo", "Munazi", "Nyagasozi", "Rukaragata", "Rwasare", "Rwigerero"}},
		{"Muhanga", "Nyabinoni", []string{"Gashorera", "Masangano", "Mbuga", "Muvumba", "Nyarusozi"}},
		{"Muhanga", "Nyamabuye", []string{"Gahogo", "Gifumba", "Gitarama", "Remera"}},
		{"Muhanga", "Nyarusange", []string{"Mbiriri", "Musongati", "Ngaru", "Rusovu"}},
		{"Muhanga", "Rongi", []string{"Gasagara", "Gasharu", "Karambo", "Nyamirambo", "Ruhango"}},
		{"Muhanga", "Rugendabari", []string{"Gasave", "Kanyana", "Kibaga", "Mpinga", "Nsanga"}},
		{"Muhanga", "Shyogwe", []string{"Kinini", "Mbare", "Mubuga", "Ruli"}},
		{"Nyamagabe", "Buruhukiro", []string{"Bushigishigi", "Byimana", "Gifurwe", "Kizimyamuriro", "Munini", "Rambya"}},
		{"Nyamagabe", "Cyanika", []string{"Gitega", "Karama", "Kiyumba", "Ngoma", "Nyanza", "Nyanzoga"}},
		{"Nyamagabe", "Gasaka", []string{"Kigeme", "Ngiryi", "Nyabivumu", "Nyamugari", "Nzega", "Remera"}},
		{"Nyamagabe", "Gatare", []string{"Bakopfu", "Gatare", "Mukongoro", "Ruganda", "Shyeru"}},
		{"Nyamagabe", "Kaduha", []string{"Kavumu", "Murambi", "Musenyi", "Nyabisindu", "Nyamiyaga"}},
		{"Nyamagabe", "Kamegeri", []string{"Bwama", "Kamegeri", "Kirehe", "Kizi", "Nyarusiza", "Rususa"}},
		{"Nyamagabe", "Kibirizi", []string{"Bugarama", "Bugarura", "Gashiha", "Karambo", "Ruhunga", "Uwindekezi"}},
		{"Nyamagabe", "Kibumbwe", []string{"Bwenda", "Gakanka", "Kibibi", "Nyakiza"}},
		{"Nyamagabe", "Kitabi", []string{"Kagano", "Mujuga", "Mukungu", "Shaba", "Uwingugu"}},
		{"Nyamagabe", "Mbazi", []string{"Manwari", "Mutiwingoma", "Ngambi", "Ngara"}},
		{"Nyamagabe", "Mugano", []string{"Gitondorero", "Gitwa", "Ruhinga", "Sovu", "Suti", "Yonde"}},
		{"Nyamagabe", "Musange", []string{"Gasave", "Jenda", "Masagara", "Masangano", "Masizi", "Nyagisozi"}},
		{"Nyamagabe", "Musebeya", []string{"Gatovu", "Nyarurambi", "Rugano", "Runege", "Rusekera", "Sekera"}},
		{"Nyamagabe", "Mushubi", []string{"Buteteri", "Cyobe", "Gashwati"}},
		{"Nyamagabe", "Nkomane", []string{"Bitandara", "Musaraba", "Mutengeri", "Nkomane", "Nyarwungo", "Twiya"}},
		{"Nyamagabe", "Tare", []string{"Buhoro", "Gasarenda", "Gatovu", "Kaganza", "Nkumbure", "Nyamigina"}},
		{"Nyamagabe", "Uwinkingi", []string{"Bigumira", "Gahira", "Kibyagira", "Mudasomwa", "Munyege", "Rugogwe"}},
		{"Nyanza", "Busasamana", []string{"Gahondo", "Kavumu", "Kibinja", "Nyanza", "Rwesero"}},
		{"Nyanza", "Busoro", []string{"Gitovu", "Kimirama", "Masangano", "Munyinya", "Rukingiro", "Shyira"}},
		{"Nyanza", "Cyabakamyi", []string{"Kadaho", "Karama", "Nyabinyenga", "Nyarurama", "Rubona"}},
		{"Nyanza", "Kibirizi", []string{"Cyeru", "Mbuye", "Mututu", "Rwotso"}},
		{"Nyanza", "Kigoma", []string{"Butansinda", "Butara", "Gahombo", "Gasoro", "Mulinja"}},
		{"Nyanza", "Mukingo", []string{"Cyerezo", "Gatagara", "Kiruli", "Mpanga", "Ngwa", "Nkomero"}},
		{"Nyanza", "Muyira", []string{"Gati", "Migina", "Nyamiyaga", "Nyamure", "Nyundo"}},
		{"Nyanza", "Ntyazo", []string{"Bugali", "Cyotamakara", "Kagunga", "Katarara"}},
		{"Nyanza", "Nyagisozi", []string{"Gahunga", "Kabirizi", "Kabuga", "Kirambi", "Rurangazi"}},
		{"Nyanza", "Rwabicuma", []string{"Gacu", "Gishike", "Mubuga", "Mushirarungu", "Nyarusange", "Runga"}},
		{"Nyaruguru", "Busanze", []string{"Kirarangombe", "Nkanda", "Nteko", "Runyombyi", "Shororo"}},
		{"Nyaruguru", "Cyahinda", []string{"Coko", "Cyahinda", "Gasasa", "Muhambara", "Rutobwe"}},
		{"Nyaruguru", "Kibeho", []string{"Gakoma", "Kibeho", "Mbasa", "Mpanda", "Mubuga", "Nyange"}},
		{"Nyaruguru", "Kivu", []string{"Cyanyirankora", "Gahurizo", "Kimina", "Kivu", "Rugerero"}},
		{"Nyaruguru", "Mata", []string{"Gorwe", "Murambi", "Nyamabuye", "Ramba", "Rwamiko"}},
		{"Nyaruguru", "Muganza", []string{"Muganza", "Rukore", "Samiyonga", "Uwacyiza"}},
		{"Nyaruguru", "Munini", []string{"Giheta", "Ngarurira", "Ngeri", "Ntwali", "Nyarure"}},
		{"Nyaruguru", "Ngera", []string{"Bitare", "Mukuge", "Murama", "Nyamirama", "Nyanza", "Yaramba"}},
		{"Nyaruguru", "Ngoma", []string{"Fugi", "Kibangu", "Kiyonza", "Mbuye", "Nyamirama", "Rubona"}},
		{"Nyaruguru", "Nyabimata", []string{"Gihemvu", "Kabere", "Mishungero", "Nyabimata", "Ruhinga"}},
		{"Nyaruguru", "Nyagisozi", []string{"Maraba", "Mwoya", "Nkakwa", "Nyagisozi"}},
		{"Nyaruguru", "Ruheru", []string{"Gitita", "Kabere", "Remera", "Ruyenzi", "Uwumusebeya"}},
		{"Nyaruguru", "Ruramba", []string{"Gabiro", "Giseke", "Nyarugano", "Rugogwe", "Ruramba"}},
		{"Nyaruguru", "Rusenge", []string{"Bunge", "Cyuna", "Gikunzi", "Mariba", "Raranzige", "Rusenge"}},
		{"Ruhango", "Bweramana", []string{"Buhanda", "Gitisi", "Murama", "Rubona", "Rwinyana"}},
		{"Ruhango", "Byimana", []string{"Kamusenyi", "Kirengeri", "Mahembe", "Mpanda", "Muhororo", "Ntenyo", "Nyakabuye"}},
		{"Ruhango", "Kabagali", []string{"Bihembe", "Karambi", "Munanira", "Remera", "Rwesero", "Rwoga"}},
		{"Ruhango", "Kinazi", []string{"Burima", "Gisali", "Kinazi", "Rubona", "Rutabo"}},
		{"Ruhango", "Kinihira", []string{"Bweramvura", "Gitinda", "Kirwa", "Muyunzwe", "Nyakogo", "Rukina"}},
		{"Ruhango", "Mbuye", []string{"Cyanza", "Gisanga", "Kabuga", "Kizibere", "Mbuye", "Mwendo", "Nyakarekare"}},
		{"Ruhango", "Mwendo", []string{"Gafunzo", "Gishweru", "Kamujisho", "Kigarama", "Kubutare", "Mutara", "Nyabibugu", "Saruheshyi"}},
		{"Ruhango", "Ntongwe", []string{"Gako", "Kareba", "Kayenzi", "Kebero", "Nyagisozi", "Nyakabungo", "Nyarurama"}},
		{"Ruhango", "Ruhango", []string{"Buhoro", "Bunyogombe", "Gikoma", "Munini", "Musamo", "Nyamagana", "Rwoga", "Tambwe"}},

		// ════════════════════════════════════════
		// WESTERN PROVINCE
		// ════════════════════════════════════════
		{"Karongi", "Bwishyura", []string{"Burunga", "Gasura", "Gitarama", "Kayenzi", "Kibuye", "Kiniha", "Nyarusazi"}},
		{"Karongi", "Gishari", []string{"Birambo", "Musasa", "Mwendo", "Rugobagoba", "Tongati"}},
		{"Karongi", "Gishyita", []string{"Buhoro", "Cyanya", "Kigarama", "Munanira", "Musasa", "Ngoma"}},
		{"Karongi", "Gitesi", []string{"Gasharu", "Gitega", "Kanunga", "Kirambo", "Munanira", "Nyamiringa", "Ruhinga", "Rwariro"}},
		{"Karongi", "Mubuga", []string{"Kagabiro", "Murangara", "Nyagatovu", "Ryaruhanga"}},
		{"Karongi", "Murambi", []string{"Mubuga", "Muhororo", "Nkoto", "Nyarunyinya", "Shyembe"}},
		{"Karongi", "Murundi", []string{"Bukiro", "Kabaya", "Kamina", "Kareba", "Nyamushishi", "Nzaratsi"}},
		{"Karongi", "Mutuntu", []string{"Byogo", "Gasharu", "Gisayura", "Kanyege", "Kinyonzwe", "Murengezo", "Rwufi"}},
		{"Karongi", "Rubengera", []string{"Bubazi", "Gacaca", "Gisanze", "Gitwa", "Kibirizi", "Mataba", "Nyarugenge", "Ruragwe"}},
		{"Karongi", "Rugabano", []string{"Gisiza", "Gitega", "Gitovu", "Kabuga", "Mubuga", "Mucyimba", "Rufungo", "Rwungo", "Tyazo"}},
		{"Karongi", "Ruganda", []string{"Biguhu", "Kabingo", "Kinyovu", "Kivumu", "Nyabikeri", "Nyamugwagwa", "Rubona", "Rugobagoba"}},
		{"Karongi", "Rwankuba", []string{"Bigugu", "Bisesero", "Gasata", "Munini", "Nyakamira", "Nyarusanga", "Rubazo", "Rubumba"}},
		{"Karongi", "Twumba", []string{"Bihumbe", "Gakuta", "Gisovu", "Gitabura", "Kavumu", "Murehe", "Rutabi"}},
		{"Ngororero", "Bwira", []string{"Bungwe", "Cyahafi", "Gashubi", "Kabarondo", "Ruhindage"}},
		{"Ngororero", "Gatumba", []string{"Cyome", "Gatsibo", "Kamasiga", "Karambo", "Ruhanga", "Rusumo"}},
		{"Ngororero", "Hindiro", []string{"Gatare", "Gatega", "Kajinge", "Marantima", "Rugendabari", "Runyinya"}},
		{"Ngororero", "Kabaya", []string{"Busunzu", "Gaseke", "Kabaya", "Mwendo", "Ngoma", "Nyenyeri"}},
		{"Ngororero", "Kageyo", []string{"Kageshi", "Kirwa", "Mukore", "Muramba", "Nyamata", "Rwamamara"}},
		{"Ngororero", "Kavumu", []string{"Birembo", "Gitwa", "Murinzi", "Nyamugeyo", "Rugeshi", "Tetero"}},
		{"Ngororero", "Matyazo", []string{"Binana", "Gitega", "Matare", "Rutare", "Rwamiko"}},
		{"Ngororero", "Muhanda", []string{"Bugarura", "Gasiza", "Mashya", "Nganzo", "Ngoma", "Rutagara"}},
		{"Ngororero", "Muhororo", []string{"Bweramana", "Mubuga", "Myiha", "Rugogwe", "Rusororo", "Sanza"}},
		{"Ngororero", "Ndaro", []string{"Bijyojyo", "Bitabage", "Kabageshi", "Kibanda", "Kinyovi"}},
		{"Ngororero", "Ngororero", []string{"Kaseke", "Kazabe", "Mugano", "Nyange", "Rususa", "Torero"}},
		{"Ngororero", "Nyange", []string{"Bambiro", "Gaseke", "Nsibo", "Vuganyana"}},
		{"Ngororero", "Sovu", []string{"Birembo", "Kagano", "Kanyana", "Musenyi", "Nyabipfura", "Rutovu"}},
		{"Nyabihu", "Bigogwe", []string{"Arusha", "Basumba", "Kijote", "Kora", "Muhe", "Rega"}},
		{"Nyabihu", "Jenda", []string{"Bukinanyana", "Gasizi", "Kabatezi", "Kareba", "Nyirakigugu", "Rega"}},
		{"Nyabihu", "Jomba", []string{"Gasiza", "Gasura", "Gisizi", "Guriro", "Kavumu", "Nyamitanzi"}},
		{"Nyabihu", "Kabatwa", []string{"Batikoti", "Cyamvumba", "Gihorwe", "Myuga", "Ngando", "Rugarama"}},
		{"Nyabihu", "Karago", []string{"Busoro", "Cyamabuye", "Gatagara", "Gihirwa", "Kadahenda", "Karengera"}},
		{"Nyabihu", "Kintobo", []string{"Gatovu", "Kintobo", "Nyagisozi", "Nyamugari", "Rukondo", "Ryinyo"}},
		{"Nyabihu", "Mukamira", []string{"Gasizi", "Jaba", "Kanyove", "Rubaya", "Rugeshi", "Rukoma", "Rurengeri"}},
		{"Nyabihu", "Muringa", []string{"Gisizi", "Mulinga", "Mwiyanike", "Nkomane", "Nyamasheke", "Rwantobo"}},
		{"Nyabihu", "Rambura", []string{"Birembo", "Guriro", "Kibisabo", "Mutaho", "Nyundo", "Rugamba"}},
		{"Nyabihu", "Rugera", []string{"Gakoro", "Marangara", "Nyagahondo", "Nyarutembe", "Rurembo", "Tyazo"}},
		{"Nyabihu", "Rurembo", []string{"Gahondo", "Gitega", "Kirimbogo", "Murambi", "Mwana", "Rwaza"}},
		{"Nyabihu", "Shyira", []string{"Cyimanzovu", "Kanyamitana", "Kintarure", "Mpinga", "Mutanda", "Shaki"}},
		{"Nyamasheke", "Bushekeri", []string{"Buvungira", "Mpumbu", "Ngoma", "Nyarusange"}},
		{"Nyamasheke", "Bushenge", []string{"Gasheke", "Impala", "Kagatamu", "Karusimbi"}},
		{"Nyamasheke", "Cyato", []string{"Bisumo", "Murambi", "Mutongo", "Rugari"}},
		{"Nyamasheke", "Gihombo", []string{"Butare", "Gitwa", "Jarama", "Kibingo", "Mubuga"}},
		{"Nyamasheke", "Kagano", []string{"Gako", "Mubumbano", "Ninzi", "Rwesero", "Shara"}},
		{"Nyamasheke", "Kanjongo", []string{"Kibogora", "Kigarama", "Kigoya", "Raro", "Susa"}},
		{"Nyamasheke", "Karambi", []string{"Gasovu", "Gitwe", "Kabuga", "Kagarama", "Rushyarara"}},
		{"Nyamasheke", "Karengera", []string{"Gasayo", "Gashashi", "Higiro", "Miko", "Mwezi"}},
		{"Nyamasheke", "Kirimbi", []string{"Cyimpindu", "Karengera", "Muhororo", "Nyarusange"}},
		{"Nyamasheke", "Macuba", []string{"Gatare", "Mutongo", "Nyakabingo", "Rugari", "Vugangoma"}},
		{"Nyamasheke", "Mahembe", []string{"Gisoke", "Kagarama", "Nyagatare", "Nyakavumu"}},
		{"Nyamasheke", "Nyabitekeri", []string{"Kigabiro", "Kinunga", "Mariba", "Muyange", "Ntango"}},
		{"Nyamasheke", "Rangiro", []string{"Banda", "Gakenke", "Jurwe", "Murambi"}},
		{"Nyamasheke", "Ruharambuga", []string{"Kanazi", "Ntendezi", "Save", "Wimana"}},
		{"Nyamasheke", "Shangi", []string{"Burimba", "Mataba", "Mugera", "Nyamugari", "Shangi"}},
		{"Rubavu", "Bugeshi", []string{"Buringo", "Butaka", "Hehu", "Kabumba", "Mutovu", "Nsherima", "Rusiza"}},
		{"Rubavu", "Busasamana", []string{"Gacurabwenge", "Gasiza", "Gihonga", "Kageshi", "Makoro", "Nyacyonga", "Rusura"}},
		{"Rubavu", "Cyanzarwe", []string{"Busigari", "Cyanzarwe", "Gora", "Kinyanzovu", "Makurizo", "Rwangara", "Rwanzekuma", "Ryabizige"}},
		{"Rubavu", "Gisenyi", []string{"Amahoro", "Bugoyi", "Kivumu", "Mbugangari", "Nengo", "Rubavu", "Umuganda"}},
		{"Rubavu", "Kanama", []string{"Kamuhoza", "Karambo", "Mahoko", "Musabike", "Nkomane", "Rusongati", "Yungwe"}},
		{"Rubavu", "Kanzenze", []string{"Kanyirabigogo", "Kirerema", "Muramba", "Nyamikongi", "Nyamirango", "Nyaruteme"}},
		{"Rubavu", "Mudende", []string{"Bihungwe", "Kanyundo", "Micinyiro", "Mirindi", "Ndururanyi", "Rungu", "Rwanyakayaga"}},
		{"Rubavu", "Nyakiriba", []string{"Bisizi", "Gikombe", "Kanyefurwe", "Nyarushyamba"}},
		{"Rubavu", "Nyamyumba", []string{"Burushya", "Busoro", "Kinigi", "Kiraga", "Munanira", "Rubona"}},
		{"Rubavu", "Nyundo", []string{"Bahimba", "Gatovu", "Kavomo", "Kigarama", "Mukondo", "Nyundo", "Terimbere"}},
		{"Rubavu", "Rubavu", []string{"Buhaza", "Burinda", "Byahi", "Gikombe", "Murambi", "Murara", "Rukoko"}},
		{"Rubavu", "Rugerero", []string{"Basa", "Gisa", "Kabilizi", "Muhira", "Rugerero", "Rushubi", "Rwaza"}},
		{"Rusizi", "Bugarama", []string{"Nyange", "Pera", "Ryankana"}},
		{"Rusizi", "Butare", []string{"Butanda", "Gatereri", "Nyamihanda", "Rwambogo"}},
		{"Rusizi", "Bweyeye", []string{"Gikungu", "Kiyabo", "Murwa", "Nyamuzi", "Rasano"}},
		{"Rusizi", "Gashonga", []string{"Birembo", "Buhokoro", "Kabakobwa", "Kacyuma", "Kamurehe", "Karemereye", "Muti", "Rusayo"}},
		{"Rusizi", "Giheke", []string{"Cyendajuru", "Gakomeye", "Giheke", "Kamashangi", "Kigenge", "Ntura", "Rwega", "Turambi"}},
		{"Rusizi", "Gihundwe", []string{"Burunga", "Gatsiro", "Gihaya", "Kagara", "Kamatita", "Shagasha"}},
		{"Rusizi", "Gikundamvura", []string{"Kizura", "Mpinga", "Nyamigina"}},
		{"Rusizi", "Gitambi", []string{"Cyingwa", "Gahungeri", "Hangabashi", "Mashesha"}},
		{"Rusizi", "Kamembe", []string{"Cyangugu", "Gihundwe", "Kamashangi", "Kamurera", "Ruganda"}},
		{"Rusizi", "Muganza", []string{"Cyarukara", "Gakoni", "Shara"}},
		{"Rusizi", "Mururu", []string{"Gahinga", "Kabahinda", "Kabasigirira", "Kagarama", "Karambi", "Miko", "Tara"}},
		{"Rusizi", "Nkanka", []string{"Gitwa", "Kamanyenga", "Kangazi", "Kinyaga", "Rugabano"}},
		{"Rusizi", "Nkombo", []string{"Bigoga", "Bugarura", "Ishywa", "Kamagimbo", "Rwenje"}},
		{"Rusizi", "Nkungu", []string{"Gatare", "Kiziguro", "Mataba", "Ryamuhirwa"}},
		{"Rusizi", "Nyakabuye", []string{"Gasebeya", "Gaseke", "Kamanu", "Kiziho", "Mashyuza", "Nyabintare"}},
		{"Rusizi", "Nyakarenzo", []string{"Gatare", "Kabagina", "Kabuye", "Kanoga", "Karangiro", "Murambi", "Rusambu"}},
		{"Rusizi", "Nzahaha", []string{"Butambamo", "Kigenge", "Murya", "Nyenji", "Rebero", "Rwinzuki"}},
		{"Rusizi", "Rwimbogo", []string{"Karenge", "Muhehwe", "Mushaka", "Rubugu", "Ruganda"}},
		{"Rutsiro", "Boneza", []string{"Bushaka", "Kabihogo", "Nkira", "Remera"}},
		{"Rutsiro", "Gihango", []string{"Bugina", "CongoNil", "Mataba", "Murambi", "Ruhingo", "Shyembe", "Teba"}},
		{"Rutsiro", "Kigeyo", []string{"Buhindure", "Nkora", "Nyagahinika", "Rukaragata"}},
		{"Rutsiro", "Kivumu", []string{"Bunyoni", "Bunyunju", "Kabere", "Kabujenje", "Karambi", "Nganzo"}},
		{"Rutsiro", "Manihira", []string{"Haniro", "Muyira", "Tangabo"}},
		{"Rutsiro", "Mukura", []string{"Kabuga", "Kagano", "Kageyo", "Kagusa", "Karambo", "Mwendo"}},
		{"Rutsiro", "Murunda", []string{"Kirwa", "Mburamazi", "Rugeyo", "Twabugezi"}},
		{"Rutsiro", "Musasa", []string{"Gabiro", "Gisiza", "Murambi", "Nyarubuye"}},
		{"Rutsiro", "Mushonyi", []string{"Biruyi", "Kaguriro", "Magaba", "Rurara"}},
		{"Rutsiro", "Mushubati", []string{"Bumba", "Cyarusera", "Gitwa", "Mageragere", "Sure"}},
		{"Rutsiro", "Nyabirasi", []string{"Busuku", "Cyivugiza", "Mubuga", "Ngoma", "Terimbere"}},
		{"Rutsiro", "Ruhango", []string{"Gatare", "Gihira", "Kavumu", "Nyakarera", "Rugasa", "Rundoyi"}},
		{"Rutsiro", "Rusebeya", []string{"Kabona", "Mberi", "Remera", "Ruronde"}},

	}

	for _, d := range data {
		var sid string
		err := tx.QueryRow(
			"SELECT s.id FROM sectors s JOIN districts d ON s.district_id = d.id WHERE s.name = ? AND d.name = ?",
			d.sector, d.district,
		).Scan(&sid)
		if err != nil {
			return err
		}
		for _, c := range d.cells {
			if _, err := tx.Exec("INSERT INTO cells (id, sector_id, name) VALUES (?, ?, ?)", newID(), sid, c); err != nil {
				return err
			}
		}
	}
	return nil
}

func seedVillages(tx *sql.Tx) error {
	// Seed a representative sample of villages for Kimironko sector cells (Gasabo district)
	type villageData struct {
		district string
		sector   string
		cell     string
		villages []string
	}

	data := []villageData{
		{"Gasabo", "Kimironko", "Kibagabaga", []string{"Ingenzi", "Umucyo", "Urumuri", "Ikaze"}},
		{"Gasabo", "Kimironko", "Bibare", []string{"Intsinzi", "Izuba", "Urukundo"}},
		{"Gasabo", "Kimironko", "Nyagatovu", []string{"Imena", "Isibo", "Ubumwe"}},
	}

	for _, d := range data {
		var cid string
		err := tx.QueryRow(
			`SELECT c.id FROM cells c
			 JOIN sectors s ON c.sector_id = s.id
			 JOIN districts d ON s.district_id = d.id
			 WHERE c.name = ? AND s.name = ? AND d.name = ?`,
			d.cell, d.sector, d.district,
		).Scan(&cid)
		if err != nil {
			continue
		}
		for _, v := range d.villages {
			if _, err := tx.Exec("INSERT INTO villages (id, cell_id, name) VALUES (?, ?, ?)", newID(), cid, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func seedBanks(tx *sql.Tx) error {
	type branchInfo struct {
		Name, Location, Phone string
	}
	type bankData struct {
		Name, ShortName, SwiftCode, Phone string
		Branches                          []branchInfo
	}

	banks := []bankData{
		{
			Name: "National Bank of Rwanda", ShortName: "BNR", SwiftCode: "ABORWRWR", Phone: "+250 252 574 282",
			Branches: []branchInfo{
				{"BNR Headquarters", "Kigali, KN 6 Ave 4", "+250 252 574 282"},
				{"BNR Rubavu Branch", "Rubavu, Western Province", ""},
			},
		},
		{
			Name: "Bank of Kigali", ShortName: "BK", SwiftCode: "BKIGRWRW", Phone: "+250 788 143 000",
			Branches: []branchInfo{
				{"BK Headquarters", "Kigali, KN 4 Ave, Plot 6112", "+250 252 593 100"},
				{"BK Kacyiru", "Kigali, Kacyiru", "+250 788 302 461"},
				{"BK Remera", "Kigali, Remera", "+250 788 304 957"},
				{"BK Remera 2", "Kigali, Remera", "+250 788 305 411"},
				{"BK Kimironko", "Kigali, Kimironko", "+250 788 305 508"},
				{"BK Manor", "Kigali, Nyarugenge", "+250 252 593 103"},
				{"BK RDB", "Kigali, Kimihurura", "+250 788 640 281"},
				{"BK KCT", "Kigali, Nyarugenge", "+250 788 640 281"},
				{"BK Kigali Market", "Kigali, Nyarugenge", "+250 788 464 15"},
				{"BK Airport", "Kigali, Kanombe", "+250 788 305 163"},
				{"BK Kicukiro", "Kigali, Kicukiro", "+250 788 386 580"},
				{"BK SFB", "Kigali, Nyarugenge", "+250 788 381 569"},
				{"BK Nyamirambo", "Kigali, Nyamirambo", "+250 788 387 022"},
				{"BK Nyabugogo", "Kigali, Nyabugogo", "+250 788 302 472"},
				{"BK Town Branch", "Kigali, Nyarugenge", "+250 788 302 514"},
				{"BK Kabuga", "Kigali, Kabuga", "+250 788 301 215"},
				{"BK Gatsata", "Kigali, Gatsata", "+250 788 302 163"},
				{"BK Rwamagana", "Rwamagana, Eastern Province", "+250 788 302 471"},
				{"BK Kayonza", "Kayonza, Eastern Province", "+250 788 301 214"},
				{"BK Ngoma", "Ngoma, Eastern Province", "+250 788 389 712"},
				{"BK Nyagatare", "Nyagatare, Eastern Province", "+250 788 383 952"},
				{"BK Ruhuha", "Bugesera, Eastern Province", "+250 788 389 897"},
				{"BK Kabarore", "Gatsibo, Eastern Province", "+250 788 389 706"},
				{"BK Gatuna", "Gicumbi, Northern Province", "+250 788 309 013"},
				{"BK Nyamata", "Bugesera, Eastern Province", "+250 788 303 995"},
				{"BK Byangabo", "Musanze, Northern Province", "+250 788 389 881"},
				{"BK Gakenke", "Gakenke, Northern Province", "+250 788 389 704"},
				{"BK Rubavu", "Rubavu, Western Province", "+250 788 302 068"},
				{"BK Huye", "Huye, Southern Province", "+250 788 302 068"},
				{"BK Muhanga", "Muhanga, Southern Province", ""},
				{"BK Rusizi", "Rusizi, Western Province", ""},
			},
		},
		{
			Name: "BPR Bank Rwanda", ShortName: "BPR", SwiftCode: "BPRBRWRW", Phone: "+250 788 187 200",
			Branches: []branchInfo{
				{"BPR Headquarters", "Kigali, KN 67 St", "+250 788 187 200"},
				{"BPR Kicukiro", "Kigali, Kicukiro", ""},
				{"BPR Remera", "Kigali, Remera", ""},
				{"BPR Nyabugogo", "Kigali, Nyabugogo", ""},
				{"BPR Kimironko", "Kigali, Kimironko", ""},
				{"BPR Nyamirambo", "Kigali, Nyamirambo", ""},
				{"BPR Rubavu", "Rubavu, Western Province", ""},
				{"BPR Musanze", "Musanze, Northern Province", ""},
				{"BPR Huye", "Huye, Southern Province", ""},
				{"BPR Muhanga", "Muhanga, Southern Province", ""},
				{"BPR Rusizi", "Rusizi, Western Province", ""},
				{"BPR Nyagatare", "Nyagatare, Eastern Province", ""},
				{"BPR Rwamagana", "Rwamagana, Eastern Province", ""},
				{"BPR Kayonza", "Kayonza, Eastern Province", ""},
				{"BPR Karongi", "Karongi, Western Province", ""},
				{"BPR Nyanza", "Nyanza, Southern Province", ""},
				{"BPR Gicumbi", "Gicumbi, Northern Province", ""},
				{"BPR Ngoma", "Ngoma, Eastern Province", ""},
				{"BPR Bugesera", "Bugesera, Eastern Province", ""},
				{"BPR Ruhango", "Ruhango, Southern Province", ""},
				{"BPR Gakenke", "Gakenke, Northern Province", ""},
			},
		},
		{
			Name: "Equity Bank Rwanda", ShortName: "Equity", SwiftCode: "EABORWRW", Phone: "+250 788 190 000",
			Branches: []branchInfo{
				{"Equity Headquarters", "Kigali, Grand Pension Plaza, Nyarugenge", "+250 788 190 000"},
				{"Equity Kisimenti", "Kigali, KG 11 Ave, Remera", ""},
				{"Equity Remera", "Kigali, Remera", ""},
				{"Equity Kicukiro", "Kigali, Kicukiro", ""},
				{"Equity Nyabugogo", "Kigali, Nyabugogo", ""},
				{"Equity Kimironko", "Kigali, Kimironko", ""},
				{"Equity Gisozi", "Kigali, Gisozi", ""},
				{"Equity Rubavu", "Rubavu, Western Province", ""},
				{"Equity Musanze", "Musanze, Northern Province", ""},
				{"Equity Rusizi/Kamembe", "Rusizi, Western Province", ""},
				{"Equity Huye", "Huye, Southern Province", ""},
				{"Equity Muhanga", "Muhanga, Southern Province", ""},
				{"Equity Nyagatare", "Nyagatare, Eastern Province", ""},
				{"Equity Kabarondo", "Kayonza, Eastern Province", ""},
				{"Equity Nyamagabe", "Nyamagabe, Southern Province", ""},
				{"Equity Nyamata", "Bugesera, Eastern Province", ""},
				{"Equity Rwamagana", "Rwamagana, Eastern Province", ""},
				{"Equity Karongi", "Karongi, Western Province", ""},
			},
		},
		{
			Name: "I&M Bank Rwanda", ShortName: "I&M", SwiftCode: "IMABORWRW", Phone: "+250 788 162 006",
			Branches: []branchInfo{
				{"I&M Head Office", "Kigali, KN 3 Ave", "+250 788 162 006"},
				{"I&M Kacyiru", "Kigali, Kacyiru", ""},
				{"I&M Remera", "Kigali, Remera", ""},
				{"I&M Kimironko", "Kigali, Kimironko", ""},
				{"I&M Nyabugogo", "Kigali, Nyabugogo", ""},
				{"I&M Kicukiro", "Kigali, Kicukiro", ""},
				{"I&M Musanze", "Musanze, Northern Province", ""},
				{"I&M Rubavu", "Rubavu, Western Province", ""},
				{"I&M Huye", "Huye, Southern Province", ""},
				{"I&M Rusizi", "Rusizi, Western Province", ""},
				{"I&M Muhanga", "Muhanga, Southern Province", ""},
			},
		},
		{
			Name: "KCB Bank Rwanda", ShortName: "KCB", SwiftCode: "KCBLRWRW", Phone: "+250 788 140 000",
			Branches: []branchInfo{
				{"KCB Headquarters", "Kigali, KN 63 St", "+250 788 140 000"},
				{"KCB Remera", "Kigali, Remera", ""},
				{"KCB Muhanga", "Muhanga, Southern Province", ""},
				{"KCB Huye", "Huye, Southern Province", ""},
				{"KCB Rusizi", "Rusizi, Western Province", ""},
				{"KCB Musanze", "Musanze, Northern Province", ""},
				{"KCB Rubavu", "Rubavu, Western Province", ""},
				{"KCB Kayonza", "Kayonza, Eastern Province", ""},
			},
		},
		{
			Name: "Guaranty Trust Bank Rwanda", ShortName: "GTBank", SwiftCode: "GTBIRWRK", Phone: "+250 788 149 600",
			Branches: []branchInfo{
				{"GTBank Head Office", "Kigali, KN 2 Ave, MIC Building", "+250 788 149 600"},
				{"GTBank Main Branch", "Kigali, KG 11 Ave, Kisimenti", ""},
				{"GTBank Gisozi", "Kigali, Gisozi", ""},
				{"GTBank Kicukiro", "Kigali, Kicukiro", ""},
				{"GTBank Kimironko", "Kigali, Kimironko", ""},
				{"GTBank Kigali City Market", "Kigali, Nyarugenge", ""},
				{"GTBank Nyabugogo", "Kigali, Nyabugogo", ""},
				{"GTBank Remera", "Kigali, Remera", ""},
				{"GTBank Karongi", "Karongi, Western Province", ""},
				{"GTBank Kayonza", "Kayonza, Eastern Province", ""},
				{"GTBank Muhanga", "Muhanga, Southern Province", ""},
				{"GTBank Musanze", "Musanze, Northern Province", ""},
				{"GTBank Ngoma", "Ngoma, Eastern Province", ""},
				{"GTBank Rubavu", "Rubavu, Western Province", ""},
				{"GTBank Rusizi", "Rusizi, Western Province", ""},
			},
		},
		{
			Name: "Access Bank Rwanda", ShortName: "Access", SwiftCode: "ABORWRWX", Phone: "+250 788 145 300",
			Branches: []branchInfo{
				{"Access Bank Head Office", "Kigali, KN 4 Ave, KIC Building", "+250 788 145 300"},
				{"Access Bank Remera", "Kigali, Remera", ""},
				{"Access Bank Nyabugogo", "Kigali, Nyabugogo", ""},
				{"Access Bank Kicukiro", "Kigali, Kicukiro", ""},
				{"Access Bank CHIC", "Kigali, Kimihurura", ""},
				{"Access Bank Rubavu", "Rubavu, Western Province", ""},
				{"Access Bank Musanze", "Musanze, Northern Province", ""},
				{"Access Bank Rusizi", "Rusizi, Western Province", ""},
			},
		},
		{
			Name: "Ecobank Rwanda", ShortName: "Ecobank", SwiftCode: "ECABORWRW", Phone: "+250 788 161 000",
			Branches: []branchInfo{
				{"Ecobank Head Office", "Kigali, KN 3 Ave, Plot 314", "+250 788 161 000"},
				{"Ecobank City Plaza", "Kigali, Avenue Commerciale", ""},
				{"Ecobank KBC", "Kigali, Kigali Business Center", ""},
				{"Ecobank Kimihurura", "Kigali, Kimihurura", ""},
				{"Ecobank Remera", "Kigali, Remera", ""},
				{"Ecobank Kimironko", "Kigali, Kimironko", ""},
				{"Ecobank Rubavu", "Rubavu, Western Province", ""},
				{"Ecobank Rusizi", "Rusizi, Western Province", ""},
				{"Ecobank Huye", "Huye, Southern Province", ""},
				{"Ecobank Muhanga", "Muhanga, Southern Province", ""},
				{"Ecobank Nyagatare", "Nyagatare, Eastern Province", ""},
				{"Ecobank Nyamagabe", "Nyamagabe, Southern Province", ""},
			},
		},
		{
			Name: "NCBA Bank Rwanda", ShortName: "NCBA", SwiftCode: "CBAFKENA", Phone: "+250 252 596 800",
			Branches: []branchInfo{
				{"NCBA Head Office", "Kigali, KG 7 Ave", "+250 252 596 800"},
				{"NCBA Remera", "Kigali, Remera", ""},
				{"NCBA Rubavu", "Rubavu, Western Province", ""},
			},
		},
	}

	for _, b := range banks {
		bankID := newID()
		if _, err := tx.Exec("INSERT INTO banks (id, name, short_name, swift_code, phone) VALUES (?, ?, ?, ?, ?)",
			bankID, b.Name, b.ShortName, b.SwiftCode, b.Phone); err != nil {
			return err
		}
		for _, br := range b.Branches {
			if _, err := tx.Exec("INSERT INTO branches (id, bank_id, name, location, phone) VALUES (?, ?, ?, ?, ?)",
				newID(), bankID, br.Name, br.Location, br.Phone); err != nil {
				return err
			}
		}
	}
	return nil
}

// seedExchangeRates populates static BNR exchange rates as a fallback.
// TODO: Connect to BNR (National Bank of Rwanda) live data feed to fetch real-time rates.
// BNR publishes daily exchange rates at https://www.bnr.rw/currency/ — consider scraping
// or using their data feed if one becomes available.
func seedExchangeRates(tx *sql.Tx) error {
	var count int
	if err := tx.QueryRow("SELECT COUNT(*) FROM exchange_rates").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	type rateData struct {
		Code    string
		Name    string
		Buying  float64
		Selling float64
		Middle  float64
	}

	// Realistic rates as of 2024-2025 (RWF per 1 unit of foreign currency)
	rates := []rateData{
		{"USD", "US Dollar", 1295.00, 1315.00, 1305.00},
		{"EUR", "Euro", 1390.00, 1420.00, 1405.00},
		{"GBP", "British Pound", 1630.00, 1670.00, 1650.00},
		{"JPY", "Japanese Yen", 8.50, 8.80, 8.65},
		{"CAD", "Canadian Dollar", 950.00, 975.00, 962.50},
		{"CHF", "Swiss Franc", 1460.00, 1500.00, 1480.00},
		{"AUD", "Australian Dollar", 850.00, 875.00, 862.50},
		{"CNY", "Chinese Yuan", 179.00, 184.00, 181.50},
		{"INR", "Indian Rupee", 15.50, 16.00, 15.75},
		{"ZAR", "South African Rand", 70.00, 73.00, 71.50},
		{"KES", "Kenyan Shilling", 9.80, 10.20, 10.00},
		{"UGX", "Ugandan Shilling", 0.34, 0.36, 0.35},
		{"TZS", "Tanzanian Shilling", 0.49, 0.52, 0.505},
		{"BIF", "Burundian Franc", 0.44, 0.47, 0.455},
		{"CDF", "Congolese Franc", 0.46, 0.49, 0.475},
		{"AED", "UAE Dirham", 352.00, 360.00, 356.00},
		{"SAR", "Saudi Riyal", 345.00, 352.00, 348.50},
		{"SDR", "Special Drawing Rights (IMF)", 1720.00, 1760.00, 1740.00},
	}

	now := "2025-01-15T10:00:00Z"
	for _, rate := range rates {
		if _, err := tx.Exec(
			"INSERT INTO exchange_rates (id, currency_code, currency_name, buying, selling, middle, date, fetched_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			newID(), rate.Code, rate.Name, rate.Buying, rate.Selling, rate.Middle, "2025-01-15", now,
		); err != nil {
			return err
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
			"INSERT INTO tariffs (id, utility, category, description, min_usage, max_usage, unit, price_per_unit, currency) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'RWF')",
			newID(), t.Utility, t.Category, t.Description, t.MinUsage, t.MaxUsage, t.Unit, t.PricePerUnit,
		); err != nil {
			return err
		}
	}
	return nil
}

func seedTelecoms(tx *sql.Tx) error {
	type planData struct {
		Name, Type, Speed, DataCap, Price, Validity string
	}
	type telecomData struct {
		Name, ShortName, Type, Website, CustomerCare string
		USSDCodes                                    []string
		Plans                                        []planData
	}

	telecoms := []telecomData{
		{
			Name: "MTN Rwanda", ShortName: "MTN", Type: "mobile", Website: "https://www.mtn.co.rw", CustomerCare: "456",
			USSDCodes: []string{"*182#:Main Menu", "*131#:Check Balance", "*345#:MTN MoMo", "*131*1#:Buy Bundles"},
			Plans: []planData{
				{"Daily 100MB", "data", "", "100MB", "200", "1 day"},
				{"Daily 500MB", "data", "", "500MB", "500", "1 day"},
				{"Weekly 1GB", "data", "", "1GB", "1500", "7 days"},
				{"Monthly 3GB", "data", "", "3GB", "5000", "30 days"},
				{"Monthly 10GB", "data", "", "10GB", "12000", "30 days"},
				{"Monthly 25GB", "data", "", "25GB", "25000", "30 days"},
				{"Tuzamurane 200", "voice", "", "", "200", "1 day"},
				{"Tuzamurane 500", "voice", "", "", "500", "3 days"},
				{"Tuzamurane 1000", "voice", "", "", "1000", "7 days"},
			},
		},
		{
			Name: "Airtel Rwanda", ShortName: "Airtel", Type: "mobile", Website: "https://www.airtel.co.rw", CustomerCare: "150",
			USSDCodes: []string{"*130#:Main Menu", "*131#:Check Balance", "*185#:Airtel Money", "*544#:Buy Bundles"},
			Plans: []planData{
				{"Daily 100MB", "data", "", "100MB", "200", "1 day"},
				{"Daily 500MB", "data", "", "500MB", "500", "1 day"},
				{"Weekly 1GB", "data", "", "1GB", "1500", "7 days"},
				{"Monthly 3GB", "data", "", "3GB", "5000", "30 days"},
				{"Monthly 10GB", "data", "", "10GB", "12000", "30 days"},
				{"Monthly 30GB", "data", "", "30GB", "25000", "30 days"},
			},
		},
		{
			Name: "BSC (Broadband Systems Corporation)", ShortName: "BSC", Type: "isp", Website: "https://www.bsc.rw", CustomerCare: "+250 252 280 888",
			Plans: []planData{
				{"Home Basic", "broadband", "10 Mbps", "Unlimited", "25000", "monthly"},
				{"Home Standard", "broadband", "25 Mbps", "Unlimited", "45000", "monthly"},
				{"Home Premium", "broadband", "50 Mbps", "Unlimited", "75000", "monthly"},
				{"Business Basic", "broadband", "50 Mbps", "Unlimited", "120000", "monthly"},
				{"Business Premium", "broadband", "100 Mbps", "Unlimited", "250000", "monthly"},
			},
		},
		{
			Name: "Liquid Intelligent Technologies", ShortName: "Liquid", Type: "isp", Website: "https://www.liquidtelecom.com/rw", CustomerCare: "+250 788 314 314",
			Plans: []planData{
				{"Home Fibre 10", "fibre", "10 Mbps", "Unlimited", "30000", "monthly"},
				{"Home Fibre 25", "fibre", "25 Mbps", "Unlimited", "50000", "monthly"},
				{"Home Fibre 50", "fibre", "50 Mbps", "Unlimited", "80000", "monthly"},
				{"Business Fibre 100", "fibre", "100 Mbps", "Unlimited", "200000", "monthly"},
			},
		},
		{
			Name: "Canal Box (Canal+)", ShortName: "Canal Box", Type: "isp", Website: "https://www.canalbox.rw", CustomerCare: "+250 788 185 555",
			Plans: []planData{
				{"Start", "fibre", "10 Mbps", "Unlimited", "18000", "monthly"},
				{"Smart", "fibre", "20 Mbps", "Unlimited", "25000", "monthly"},
				{"Power", "fibre", "50 Mbps", "Unlimited", "45000", "monthly"},
				{"Ultra", "fibre", "100 Mbps", "Unlimited", "80000", "monthly"},
			},
		},
	}

	for _, t := range telecoms {
		telecomID := newID()
		ussdStr := ""
		for i, u := range t.USSDCodes {
			if i > 0 {
				ussdStr += ","
			}
			ussdStr += u
		}
		if _, err := tx.Exec(
			"INSERT INTO telecoms (id, name, short_name, type, website, customer_care, ussd_codes) VALUES (?, ?, ?, ?, ?, ?, ?)",
			telecomID, t.Name, t.ShortName, t.Type, t.Website, t.CustomerCare, ussdStr,
		); err != nil {
			return err
		}
		for _, p := range t.Plans {
			if _, err := tx.Exec(
				"INSERT INTO telecom_plans (id, telecom_id, name, type, speed, data_cap, price, validity) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
				newID(), telecomID, p.Name, p.Type, p.Speed, p.DataCap, p.Price, p.Validity,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func seedMobileMoney(tx *sql.Tx) error {
	type feeData struct {
		TxType, MinAmount, MaxAmount, Fee string
	}
	type serviceData struct {
		Name, Provider, USSDCode, AgentUSSD string
		DailyLimit, MonthlyLimit            string
		Fees                                []feeData
	}

	services := []serviceData{
		{
			Name: "MTN Mobile Money (MoMo)", Provider: "MTN Rwanda",
			USSDCode: "*182*8#", AgentUSSD: "*345#",
			DailyLimit: "3000000", MonthlyLimit: "20000000",
			Fees: []feeData{
				{"transfer", "1", "5000", "50"},
				{"transfer", "5001", "10000", "100"},
				{"transfer", "10001", "100000", "250"},
				{"transfer", "100001", "500000", "500"},
				{"transfer", "500001", "1000000", "1000"},
				{"transfer", "1000001", "3000000", "1500"},
				{"withdraw", "1", "5000", "300"},
				{"withdraw", "5001", "10000", "400"},
				{"withdraw", "10001", "100000", "750"},
				{"withdraw", "100001", "500000", "1500"},
				{"withdraw", "500001", "1000000", "2500"},
			},
		},
		{
			Name: "Airtel Money", Provider: "Airtel Rwanda",
			USSDCode: "*185#", AgentUSSD: "*185*9#",
			DailyLimit: "3000000", MonthlyLimit: "20000000",
			Fees: []feeData{
				{"transfer", "1", "5000", "50"},
				{"transfer", "5001", "10000", "100"},
				{"transfer", "10001", "100000", "250"},
				{"transfer", "100001", "500000", "500"},
				{"transfer", "500001", "1000000", "1000"},
				{"withdraw", "1", "5000", "300"},
				{"withdraw", "5001", "10000", "400"},
				{"withdraw", "10001", "100000", "700"},
				{"withdraw", "100001", "500000", "1400"},
				{"withdraw", "500001", "1000000", "2400"},
			},
		},
	}

	for _, s := range services {
		serviceID := newID()
		if _, err := tx.Exec(
			"INSERT INTO mobile_money (id, name, provider, ussd_code, agent_ussd, daily_limit, monthly_limit) VALUES (?, ?, ?, ?, ?, ?, ?)",
			serviceID, s.Name, s.Provider, s.USSDCode, s.AgentUSSD, s.DailyLimit, s.MonthlyLimit,
		); err != nil {
			return err
		}
		for _, f := range s.Fees {
			if _, err := tx.Exec(
				"INSERT INTO mobile_money_fees (id, service_id, tx_type, min_amount, max_amount, fee) VALUES (?, ?, ?, ?, ?, ?)",
				newID(), serviceID, f.TxType, f.MinAmount, f.MaxAmount, f.Fee,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func seedEmergencyNumbers(tx *sql.Tx) error {
	type numberData struct {
		Name, Number, Category, Description string
	}

	numbers := []numberData{
		{"Rwanda National Police", "112", "emergency", "Police emergency line"},
		{"Rwanda Investigation Bureau (RIB)", "166", "emergency", "Report crimes and criminal activity"},
		{"Ambulance / SAMU", "912", "emergency", "Medical emergency and ambulance service"},
		{"Fire Brigade", "111", "emergency", "Fire emergency response"},
		{"Gender-Based Violence Hotline", "3512", "emergency", "Report GBV cases, free 24/7 line"},
		{"Child Helpline (Turi Hano)", "116", "emergency", "Child protection and support, free 24/7"},
		{"Anti-Corruption Hotline", "997", "government", "Report corruption (Office of the Ombudsman)"},
		{"Traffic Police", "113", "emergency", "Traffic accidents and road emergencies"},
		{"Rwanda Revenue Authority (RRA)", "3004", "government", "Tax inquiries and support"},
		{"RURA Consumer Line", "3717", "government", "Rwanda Utilities Regulatory Authority"},
		{"IREMBO Support", "9099", "government", "E-government services support"},
		{"Isange One Stop Center", "3029", "health", "GBV survivors medical and legal support"},
		{"COVID-19 Hotline", "114", "health", "Health emergency and disease reporting"},
		{"RSSB", "+250 788 185 289", "government", "Rwanda Social Security Board"},
		{"MTN Customer Care", "456", "telecom", "MTN Rwanda customer support"},
		{"Airtel Customer Care", "150", "telecom", "Airtel Rwanda customer support"},
		{"EUCL (Electricity)", "2727", "utility", "Electricity complaints and outages"},
		{"WASAC (Water)", "+250 252 530 340", "utility", "Water supply issues"},
		{"RwandAir", "+250 788 177 000", "travel", "Rwanda national airline"},
		{"Kigali International Airport", "+250 252 585 800", "travel", "Airport information"},
	}

	for _, n := range numbers {
		if _, err := tx.Exec(
			"INSERT INTO emergency_numbers (id, name, number, category, description) VALUES (?, ?, ?, ?, ?)",
			newID(), n.Name, n.Number, n.Category, n.Description,
		); err != nil {
			return err
		}
	}
	return nil
}

func seedPublicHolidays(tx *sql.Tx) error {
	type holidayData struct {
		Name, Date, Description string
		IsMovable               bool
	}

	holidays := []holidayData{
		{"New Year's Day", "01-01", "Celebration of the new year", false},
		{"National Heroes' Day", "02-01", "Honors national heroes who contributed to Rwanda", false},
		{"Good Friday", "variable", "Christian holiday, Friday before Easter", true},
		{"Easter Monday", "variable", "Christian holiday, Monday after Easter", true},
		{"Genocide Against the Tutsi Memorial Day", "04-07", "Commemoration of the 1994 Genocide against the Tutsi. Begins a week of mourning (Kwibuka)", false},
		{"Labour Day", "05-01", "International Workers' Day", false},
		{"Eid al-Fitr", "variable", "End of Ramadan (Islamic holiday)", true},
		{"Independence Day", "07-01", "Anniversary of Rwanda's independence (1962)", false},
		{"Liberation Day", "07-04", "Celebrates the end of the 1994 Genocide, RPF victory", false},
		{"Umuganura Day", "08-01", "National Harvest Festival, celebrating first harvest", false},
		{"Assumption Day", "08-15", "Christian holiday honoring the Virgin Mary", false},
		{"Eid al-Adha", "variable", "Islamic Feast of Sacrifice", true},
		{"Christmas Day", "12-25", "Christian celebration of the birth of Jesus", false},
		{"Boxing Day", "12-26", "Day after Christmas", false},
	}

	for _, h := range holidays {
		movable := 0
		if h.IsMovable {
			movable = 1
		}
		if _, err := tx.Exec(
			"INSERT INTO public_holidays (id, name, date, description, is_movable) VALUES (?, ?, ?, ?, ?)",
			newID(), h.Name, h.Date, h.Description, movable,
		); err != nil {
			return err
		}
	}
	return nil
}
