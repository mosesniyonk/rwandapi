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
		// GASABO DISTRICT
		// ════════════════════════════════════════
		{"Gasabo", "Bumbogo", []string{"Nyagasamba", "Nyankongi", "Rubungo", "Ruhinga", "Terimbere"}},
		{"Gasabo", "Gatsata", []string{"Karuruma", "Kinyaga", "Nyabikenke", "Nyacyonga", "Nyamabuye"}},
		{"Gasabo", "Gikomero", []string{"Gasave", "Gikomero", "Kibara", "Mukindo", "Murama"}},
		{"Gasabo", "Gisozi", []string{"Amarembo", "Umubano", "Umujyi", "Urugero", "Umunyinya"}},
		{"Gasabo", "Jabana", []string{"Gasagara", "Giticyinyoni", "Jari", "Muganza", "Nyagasozi"}},
		{"Gasabo", "Jali", []string{"Kabuga", "Kamashashi", "Karama", "Rubilizi", "Rurima"}},
		{"Gasabo", "Kacyiru", []string{"Kamatamu", "Kamutwa", "Kibaza", "Rugando"}},
		{"Gasabo", "Kimihurura", []string{"Kamashashi", "Kibaza", "Kimihurura", "Rugando", "Ururembo"}},
		{"Gasabo", "Kimironko", []string{"Bibare", "Kibagabaga", "Kimironko", "Nyagatovu"}},
		{"Gasabo", "Kinyinya", []string{"Gacuriro", "Gasharu", "Kagugu", "Kamatamu", "Murama"}},
		{"Gasabo", "Ndera", []string{"Byimana", "Masoro", "Musave", "Ndamira", "Ruvuzo"}},
		{"Gasabo", "Nduba", []string{"Gasagara", "Kamahame", "Kibagabaga", "Muremure", "Shango"}},
		{"Gasabo", "Remera", []string{"Gishushu", "Nyabisindu", "Rukiri I", "Rukiri II"}},
		{"Gasabo", "Rusororo", []string{"Gasagara", "Kabuga", "Kinunga", "Masizi", "Mbandazi"}},
		{"Gasabo", "Rutunga", []string{"Bwerankori", "Cyimana", "Ngara", "Nyamugari", "Rutonde"}},

		// ════════════════════════════════════════
		// KICUKIRO DISTRICT
		// ════════════════════════════════════════
		{"Kicukiro", "Gahanga", []string{"Gahanga", "Kagasa", "Karembure", "Murinja", "Nunga"}},
		{"Kicukiro", "Gatenga", []string{"Gatenga", "Nyabisindu", "Rebero", "Rubirizi"}},
		{"Kicukiro", "Gikondo", []string{"Gikondo", "Kinunga", "Murambi", "Sovu"}},
		{"Kicukiro", "Kagarama", []string{"Kagarama", "Kanserege", "Muyange", "Rukatsa"}},
		{"Kicukiro", "Kanombe", []string{"Busanza", "Kabeza", "Karama", "Nyarurama"}},
		{"Kicukiro", "Kicukiro", []string{"Biryogo", "Gasharu", "Ngoma", "Niboye"}},
		{"Kicukiro", "Kigarama", []string{"Kigarama", "Muyange", "Nyakagunga", "Rwampara"}},
		{"Kicukiro", "Masaka", []string{"Cyimo", "Gasagara", "Mbabe", "Muhanga"}},
		{"Kicukiro", "Niboye", []string{"Gatare", "Kigali", "Niboye", "Nyabikenke"}},
		{"Kicukiro", "Nyarugunga", []string{"Kamashashi", "Nonko", "Nyarugunga", "Rwimbogo"}},

		// ════════════════════════════════════════
		// NYARUGENGE DISTRICT
		// ════════════════════════════════════════
		{"Nyarugenge", "Gitega", []string{"Akabahizi", "Akabeza", "Akamatamu", "Igikapu", "Inyarurembo"}},
		{"Nyarugenge", "Kanyinya", []string{"Gakoki", "Munini", "Nyamweru", "Rubona", "Rwampara"}},
		{"Nyarugenge", "Kigali", []string{"Kigali", "Nyabugogo"}},
		{"Nyarugenge", "Kimisagara", []string{"Kimisagara", "Kora", "Musezero", "Nyabugogo", "Overijse"}},
		{"Nyarugenge", "Mageragere", []string{"Kamukina", "Mugina", "Mweza", "Nkubili", "Nyamabuye"}},
		{"Nyarugenge", "Muhima", []string{"Amahoro", "Ihuriro", "Ingenzi", "Ubumwe"}},
		{"Nyarugenge", "Nyakabanda", []string{"Cyivugiza", "Inyarurembo", "Kabeza", "Mumena"}},
		{"Nyarugenge", "Nyamirambo", []string{"Cyivugiza", "Mumena", "Rugarama", "Katabaro"}},
		{"Nyarugenge", "Nyarugenge", []string{"Kabasengerezi", "Kagunga", "Kamuhoza", "Kivu"}},
		{"Nyarugenge", "Rwezamenyo", []string{"Agatare", "Iterambere", "Kagunga", "Urugwiro"}},

		// ════════════════════════════════════════
		// BUGESERA DISTRICT
		// ════════════════════════════════════════
		{"Bugesera", "Gashora", []string{"Gashora", "Kagomasi", "Karera", "Mumana", "Ramiro"}},
		{"Bugesera", "Juru", []string{"Gakomeye", "Juru", "Maranyundo", "Mugorore", "Nyabagendwa"}},
		{"Bugesera", "Kamabuye", []string{"Gishali", "Kabuye", "Kamabuye", "Mugorore", "Rurenge"}},
		{"Bugesera", "Mareba", []string{"Gakomeye", "Kabuye", "Mareba", "Mugorore", "Ruvumera"}},
		{"Bugesera", "Mayange", []string{"Gashikiri", "Kanzenze", "Mayange", "Murama", "Nyarubaka"}},
		{"Bugesera", "Musenyi", []string{"Gatare", "Gishubi", "Kibungo", "Musenyi", "Rweru"}},
		{"Bugesera", "Mwogo", []string{"Kabuye", "Kindama", "Mwogo", "Nyamure", "Ruhinga"}},
		{"Bugesera", "Ngeruka", []string{"Murama", "Ngeruka", "Ntarama", "Ruhuha", "Rurenge"}},
		{"Bugesera", "Ntarama", []string{"Cyugaro", "Kibungo", "Ntarama", "Nyamure", "Rilima"}},
		{"Bugesera", "Nyamata", []string{"Kanazi", "Kanzenze", "Nyamata", "Ringaninka", "Rubona"}},
		{"Bugesera", "Nyarugenge", []string{"Gashora", "Nyarugenge", "Rugunga", "Rurenge", "Rweru"}},
		{"Bugesera", "Rilima", []string{"Kabuye", "Ntarama", "Rilima", "Ruhuha", "Rweru"}},
		{"Bugesera", "Ruhuha", []string{"Gikundamvura", "Kindama", "Ruhuha", "Rweru", "Shyara"}},
		{"Bugesera", "Rweru", []string{"Gashora", "Kabuye", "Nemba", "Rweru", "Rwibiraro"}},
		{"Bugesera", "Shyara", []string{"Kamabuye", "Nyamata", "Ruhuha", "Shyara"}},

		// ════════════════════════════════════════
		// GATSIBO DISTRICT
		// ════════════════════════════════════════
		{"Gatsibo", "Gasange", []string{"Gasange", "Kabare", "Remera", "Rugarama", "Rwimbogo"}},
		{"Gatsibo", "Gatsibo", []string{"Gatsibo", "Kageyo", "Remera", "Rugarama"}},
		{"Gatsibo", "Gitoki", []string{"Gitoki", "Kiramuruzi", "Murambi", "Nyagihanga", "Rwimbogo"}},
		{"Gatsibo", "Kabarore", []string{"Gakenke", "Kabarore", "Kiziguro", "Matimba", "Rwimbogo"}},
		{"Gatsibo", "Kageyo", []string{"Gasange", "Kageyo", "Murambi", "Nyagahanga", "Remera"}},
		{"Gatsibo", "Kiramuruzi", []string{"Gahini", "Kabarore", "Kiramuruzi", "Murambi", "Nyagihanga"}},
		{"Gatsibo", "Kiziguro", []string{"Kigabiro", "Kiziguro", "Mukarange", "Nyagihanga", "Rugarama"}},
		{"Gatsibo", "Muhura", []string{"Kabare", "Muhura", "Nyagihanga", "Remera", "Rugarama"}},
		{"Gatsibo", "Murambi", []string{"Kageyo", "Murambi", "Nyagahanga", "Remera", "Rwimbogo"}},
		{"Gatsibo", "Ngarama", []string{"Gasange", "Kageyo", "Ngarama", "Nyagihanga", "Remera"}},
		{"Gatsibo", "Nyagihanga", []string{"Gatsibo", "Kabarore", "Nyagihanga", "Rugarama", "Rwimbogo"}},
		{"Gatsibo", "Remera", []string{"Gahini", "Kabare", "Mukarange", "Remera", "Rugarama"}},
		{"Gatsibo", "Rugarama", []string{"Gasange", "Murambi", "Rugarama", "Rwimbogo"}},
		{"Gatsibo", "Rwimbogo", []string{"Gatsibo", "Kageyo", "Murambi", "Rwimbogo"}},

		// ════════════════════════════════════════
		// KAYONZA DISTRICT
		// ════════════════════════════════════════
		{"Kayonza", "Gahini", []string{"Gahini", "Juru", "Kigabiro", "Remera", "Urugarama"}},
		{"Kayonza", "Kabare", []string{"Kabare", "Karama", "Murama", "Rubona", "Rwinkwavu"}},
		{"Kayonza", "Kabarondo", []string{"Gitaraga", "Kabarondo", "Ravin", "Rwesero"}},
		{"Kayonza", "Mukarange", []string{"Karama", "Mukarange", "Nyamirama", "Rubona"}},
		{"Kayonza", "Murama", []string{"Kabare", "Murama", "Nyamirama", "Rubona", "Ruramira"}},
		{"Kayonza", "Murundi", []string{"Gahini", "Kabarondo", "Murundi", "Rwinkwavu"}},
		{"Kayonza", "Mwiri", []string{"Kabare", "Musha", "Mwiri", "Ndego", "Rwinkwavu"}},
		{"Kayonza", "Ndego", []string{"Kabare", "Murama", "Ndego", "Nyamirama"}},
		{"Kayonza", "Nyamirama", []string{"Gahini", "Kabarondo", "Nyamirama", "Rubona"}},
		{"Kayonza", "Rukara", []string{"Kabarondo", "Musha", "Rukara", "Rwinkwavu"}},
		{"Kayonza", "Ruramira", []string{"Gahini", "Murama", "Ruramira", "Rwinkwavu"}},
		{"Kayonza", "Rwinkwavu", []string{"Kabarondo", "Ndego", "Rwinkwavu", "Cyanya"}},

		// ════════════════════════════════════════
		// KIREHE DISTRICT
		// ════════════════════════════════════════
		{"Kirehe", "Gahara", []string{"Gahara", "Nyamugari", "Rugarama", "Rwesero"}},
		{"Kirehe", "Gatore", []string{"Gahara", "Gatore", "Nyarubuye", "Rugarama"}},
		{"Kirehe", "Kigarama", []string{"Kigarama", "Munini", "Nyamugari", "Rugarama"}},
		{"Kirehe", "Kigina", []string{"Kigina", "Munini", "Nyamugari", "Rugarama"}},
		{"Kirehe", "Kirehe", []string{"Karenge", "Kirehe", "Mpanga", "Mushikiri"}},
		{"Kirehe", "Mahama", []string{"Mahama", "Munini", "Nyarubuye", "Rugarama"}},
		{"Kirehe", "Mpanga", []string{"Mpanga", "Mushikiri", "Nasho", "Nyamugari"}},
		{"Kirehe", "Musaza", []string{"Gahara", "Musaza", "Nyamugari", "Rugarama"}},
		{"Kirehe", "Mushikiri", []string{"Gahara", "Mushikiri", "Nasho", "Nyamugari"}},
		{"Kirehe", "Nasho", []string{"Munini", "Nasho", "Nyarubuye", "Rwesero"}},
		{"Kirehe", "Nyamugari", []string{"Gatore", "Kigarama", "Nyamugari", "Rugarama"}},
		{"Kirehe", "Nyarubuye", []string{"Nyamugari", "Nyarubuye", "Rugarama", "Rwesero"}},

		// ════════════════════════════════════════
		// NGOMA DISTRICT
		// ════════════════════════════════════════
		{"Ngoma", "Gashanda", []string{"Gashanda", "Karama", "Murama", "Remera"}},
		{"Ngoma", "Jarama", []string{"Jarama", "Kazo", "Murama", "Rurenge"}},
		{"Ngoma", "Karembo", []string{"Karembo", "Murama", "Remera", "Rurenge"}},
		{"Ngoma", "Kazo", []string{"Gashanda", "Kazo", "Murama", "Rurenge"}},
		{"Ngoma", "Kibungo", []string{"Cyasemakamba", "Kibungo", "Murama", "Rurenge"}},
		{"Ngoma", "Mugesera", []string{"Karama", "Mugesera", "Remera", "Rwimbogo"}},
		{"Ngoma", "Murama", []string{"Gashanda", "Karama", "Murama", "Remera"}},
		{"Ngoma", "Mutenderi", []string{"Jarama", "Murama", "Mutenderi", "Rurenge"}},
		{"Ngoma", "Remera", []string{"Gashanda", "Karama", "Remera", "Rurenge"}},
		{"Ngoma", "Rukira", []string{"Kazo", "Murama", "Rukira", "Rurenge"}},
		{"Ngoma", "Rukumberi", []string{"Karama", "Murama", "Rukumberi", "Rurenge"}},
		{"Ngoma", "Rurenge", []string{"Jarama", "Kazo", "Murama", "Rurenge"}},
		{"Ngoma", "Sake", []string{"Kazo", "Murama", "Rurenge", "Sake"}},
		{"Ngoma", "Zaza", []string{"Karama", "Murama", "Rurenge", "Zaza"}},

		// ════════════════════════════════════════
		// NYAGATARE DISTRICT
		// ════════════════════════════════════════
		{"Nyagatare", "Gatunda", []string{"Gatunda", "Karama", "Matimba", "Nyagatare", "Rwempasha"}},
		{"Nyagatare", "Karama", []string{"Karama", "Karangazi", "Matimba", "Rwimiyaga"}},
		{"Nyagatare", "Karangazi", []string{"Karangazi", "Matimba", "Nyagatare", "Rwempasha"}},
		{"Nyagatare", "Katabagemu", []string{"Gatunda", "Karangazi", "Katabagemu", "Matimba"}},
		{"Nyagatare", "Kiyombe", []string{"Karama", "Kiyombe", "Matimba", "Nyagatare"}},
		{"Nyagatare", "Matimba", []string{"Karangazi", "Matimba", "Nyagatare", "Tabagwe"}},
		{"Nyagatare", "Mimuri", []string{"Karama", "Mimuri", "Nyagatare", "Tabagwe"}},
		{"Nyagatare", "Mukama", []string{"Karangazi", "Mukama", "Nyagatare", "Rwempasha"}},
		{"Nyagatare", "Musheri", []string{"Karama", "Matimba", "Musheri", "Rwimiyaga"}},
		{"Nyagatare", "Nyagatare", []string{"Kiyombe", "Matimba", "Nyagatare", "Tabagwe"}},
		{"Nyagatare", "Rukomo", []string{"Karangazi", "Mukama", "Rukomo", "Rwempasha"}},
		{"Nyagatare", "Rwempasha", []string{"Matimba", "Nyagatare", "Rwempasha", "Tabagwe"}},
		{"Nyagatare", "Rwimiyaga", []string{"Karama", "Karangazi", "Matimba", "Rwimiyaga"}},
		{"Nyagatare", "Tabagwe", []string{"Matimba", "Nyagatare", "Rwempasha", "Tabagwe"}},

		// ════════════════════════════════════════
		// RWAMAGANA DISTRICT
		// ════════════════════════════════════════
		{"Rwamagana", "Fumbwe", []string{"Fumbwe", "Karenge", "Munyaga", "Rubona"}},
		{"Rwamagana", "Gahengeri", []string{"Gahengeri", "Karenge", "Musha", "Nzige"}},
		{"Rwamagana", "Gishali", []string{"Gishali", "Karenge", "Musha", "Mwulire"}},
		{"Rwamagana", "Karenge", []string{"Fumbwe", "Karenge", "Munyaga", "Nzige"}},
		{"Rwamagana", "Kigabiro", []string{"Gahengeri", "Kigabiro", "Musha", "Mwulire"}},
		{"Rwamagana", "Muhazi", []string{"Karenge", "Muhazi", "Munyiginya", "Rubona"}},
		{"Rwamagana", "Munyaga", []string{"Fumbwe", "Karenge", "Munyaga", "Rubona"}},
		{"Rwamagana", "Munyiginya", []string{"Karenge", "Munyiginya", "Nzige", "Rubona"}},
		{"Rwamagana", "Musha", []string{"Gahengeri", "Gishali", "Musha", "Nyakaliro"}},
		{"Rwamagana", "Muyumbu", []string{"Gahengeri", "Kigabiro", "Muyumbu", "Nzige"}},
		{"Rwamagana", "Mwulire", []string{"Gishali", "Kigabiro", "Musha", "Mwulire"}},
		{"Rwamagana", "Nyakaliro", []string{"Gahengeri", "Musha", "Nyakaliro", "Rubona"}},
		{"Rwamagana", "Nzige", []string{"Fumbwe", "Karenge", "Munyaga", "Nzige"}},
		{"Rwamagana", "Rubona", []string{"Karenge", "Muhazi", "Munyiginya", "Rubona"}},

		// ════════════════════════════════════════
		// BURERA DISTRICT
		// ════════════════════════════════════════
		{"Burera", "Bungwe", []string{"Bungwe", "Gahunga", "Kagogo", "Ruhunde"}},
		{"Burera", "Butaro", []string{"Butaro", "Cyanika", "Ruhunde", "Rugarama"}},
		{"Burera", "Cyanika", []string{"Butaro", "Cyanika", "Kagogo", "Ruhunde"}},
		{"Burera", "Cyeru", []string{"Cyeru", "Gahunga", "Gitovu", "Kagogo"}},
		{"Burera", "Gahunga", []string{"Bungwe", "Gahunga", "Kagogo", "Rugarama"}},
		{"Burera", "Gatebe", []string{"Gatebe", "Gitovu", "Kagogo", "Nemba"}},
		{"Burera", "Gitovu", []string{"Cyeru", "Gatebe", "Gitovu", "Kagogo"}},
		{"Burera", "Kagogo", []string{"Bungwe", "Gahunga", "Kagogo", "Ruhunde"}},
		{"Burera", "Kinoni", []string{"Kinoni", "Kinyababa", "Nemba", "Rugarama"}},
		{"Burera", "Kinyababa", []string{"Kinoni", "Kinyababa", "Nemba", "Ruhunde"}},
		{"Burera", "Kivuye", []string{"Gahunga", "Kagogo", "Kivuye", "Ruhunde"}},
		{"Burera", "Nemba", []string{"Gatebe", "Kinoni", "Nemba", "Rugarama"}},
		{"Burera", "Rugarama", []string{"Butaro", "Gahunga", "Rugarama", "Ruhunde"}},
		{"Burera", "Rugendabari", []string{"Kagogo", "Nemba", "Rugendabari", "Ruhunde"}},
		{"Burera", "Ruhunde", []string{"Bungwe", "Kagogo", "Ruhunde", "Rugarama"}},
		{"Burera", "Rusarabuye", []string{"Gahunga", "Kagogo", "Rusarabuye", "Ruhunde"}},
		{"Burera", "Rwerere", []string{"Gahunga", "Kagogo", "Ruhunde", "Rwerere"}},

		// ════════════════════════════════════════
		// GAKENKE DISTRICT
		// ════════════════════════════════════════
		{"Gakenke", "Busengo", []string{"Busengo", "Gashenyi", "Mugunga", "Rushashi"}},
		{"Gakenke", "Coko", []string{"Coko", "Janja", "Kamubuga", "Mugunga"}},
		{"Gakenke", "Cyabingo", []string{"Cyabingo", "Gashenyi", "Karambo", "Rushashi"}},
		{"Gakenke", "Gakenke", []string{"Gakenke", "Gashenyi", "Kamubuga", "Mataba"}},
		{"Gakenke", "Gashenyi", []string{"Busengo", "Gashenyi", "Muhondo", "Rushashi"}},
		{"Gakenke", "Janja", []string{"Coko", "Janja", "Kamubuga", "Mataba"}},
		{"Gakenke", "Kamubuga", []string{"Gakenke", "Janja", "Kamubuga", "Mugunga"}},
		{"Gakenke", "Karambo", []string{"Cyabingo", "Karambo", "Kivuruga", "Rushashi"}},
		{"Gakenke", "Kivuruga", []string{"Karambo", "Kivuruga", "Mataba", "Muhondo"}},
		{"Gakenke", "Mataba", []string{"Gakenke", "Janja", "Mataba", "Mugunga"}},
		{"Gakenke", "Minazi", []string{"Gashenyi", "Minazi", "Muhondo", "Rushashi"}},
		{"Gakenke", "Mugunga", []string{"Coko", "Kamubuga", "Mugunga", "Rushashi"}},
		{"Gakenke", "Muhondo", []string{"Gashenyi", "Kivuruga", "Minazi", "Muhondo"}},
		{"Gakenke", "Muyongwe", []string{"Kamubuga", "Mataba", "Mugunga", "Muyongwe"}},
		{"Gakenke", "Muzo", []string{"Gashenyi", "Mugunga", "Muhondo", "Muzo"}},
		{"Gakenke", "Nemba", []string{"Busengo", "Gashenyi", "Nemba", "Rushashi"}},
		{"Gakenke", "Ruli", []string{"Kamubuga", "Mataba", "Mugunga", "Ruli"}},
		{"Gakenke", "Rusasa", []string{"Gashenyi", "Muhondo", "Rushashi", "Rusasa"}},
		{"Gakenke", "Rushashi", []string{"Busengo", "Gashenyi", "Muhondo", "Rushashi"}},

		// ════════════════════════════════════════
		// GICUMBI DISTRICT
		// ════════════════════════════════════════
		{"Gicumbi", "Bukure", []string{"Bukure", "Bwisige", "Manyagiro", "Rutare"}},
		{"Gicumbi", "Bwisige", []string{"Bukure", "Bwisige", "Kaniga", "Rutare"}},
		{"Gicumbi", "Byumba", []string{"Byumba", "Kaniga", "Nyamiyaga", "Rutare"}},
		{"Gicumbi", "Cyumba", []string{"Cyumba", "Manyagiro", "Miyove", "Shangasha"}},
		{"Gicumbi", "Giti", []string{"Giti", "Kaniga", "Rutare", "Ruvune"}},
		{"Gicumbi", "Kaniga", []string{"Bwisige", "Kaniga", "Nyamiyaga", "Rutare"}},
		{"Gicumbi", "Manyagiro", []string{"Bukure", "Cyumba", "Manyagiro", "Shangasha"}},
		{"Gicumbi", "Miyove", []string{"Cyumba", "Miyove", "Muko", "Shangasha"}},
		{"Gicumbi", "Kageyo", []string{"Byumba", "Kageyo", "Kaniga", "Rutare"}},
		{"Gicumbi", "Mukarange", []string{"Bwisige", "Kaniga", "Mukarange", "Rutare"}},
		{"Gicumbi", "Muko", []string{"Cyumba", "Miyove", "Muko", "Shangasha"}},
		{"Gicumbi", "Mutete", []string{"Bukure", "Kaniga", "Mutete", "Rutare"}},
		{"Gicumbi", "Nyamiyaga", []string{"Byumba", "Kaniga", "Nyamiyaga", "Ruvune"}},
		{"Gicumbi", "Nyankenke", []string{"Kaniga", "Nyankenke", "Rutare", "Ruvune"}},
		{"Gicumbi", "Rubaya", []string{"Kaniga", "Rubaya", "Rutare", "Rwamiko"}},
		{"Gicumbi", "Rukomo", []string{"Bukure", "Manyagiro", "Rukomo", "Shangasha"}},
		{"Gicumbi", "Rushaki", []string{"Kaniga", "Miyove", "Rushaki", "Shangasha"}},
		{"Gicumbi", "Rutare", []string{"Byumba", "Giti", "Kaniga", "Rutare"}},
		{"Gicumbi", "Ruvune", []string{"Giti", "Kaniga", "Nyamiyaga", "Ruvune"}},
		{"Gicumbi", "Rwamiko", []string{"Kaniga", "Rubaya", "Rutare", "Rwamiko"}},
		{"Gicumbi", "Shangasha", []string{"Cyumba", "Manyagiro", "Miyove", "Shangasha"}},

		// ════════════════════════════════════════
		// MUSANZE DISTRICT
		// ════════════════════════════════════════
		{"Musanze", "Busogo", []string{"Busogo", "Cyabararika", "Gacaca", "Kavumu", "Sahara"}},
		{"Musanze", "Cyuve", []string{"Bihinga", "Cyuve", "Kabere", "Migeshi", "Nyarubande"}},
		{"Musanze", "Gacaca", []string{"Gacaca", "Kabaya", "Nyaruhonga", "Rubona", "Ruhanga"}},
		{"Musanze", "Gashaki", []string{"Gashaki", "Kabirizi", "Kavumu", "Rugerero", "Ruhanga"}},
		{"Musanze", "Gataraga", []string{"Gataraga", "Kaguhu", "Nyarubuye", "Ruhondo", "Rurengeri"}},
		{"Musanze", "Kimonyi", []string{"Gasiza", "Kabeza", "Kavumu", "Nyabageni"}},
		{"Musanze", "Kinigi", []string{"Bisoke", "Kampanga", "Kaguhu", "Kinigi", "Nyonirima"}},
		{"Musanze", "Muhoza", []string{"Cyivugiza", "Kampanga", "Kimonyi", "Mpenge", "Ruhongore"}},
		{"Musanze", "Muko", []string{"Gashaki", "Kavumu", "Muko", "Nganzo", "Rurengeri"}},
		{"Musanze", "Musanze", []string{"Cyuve", "Gacaca", "Mpenge", "Musanze", "Songa"}},
		{"Musanze", "Nkotsi", []string{"Kabatwa", "Nkotsi", "Nyamyumba", "Ruhondo", "Shingiro"}},
		{"Musanze", "Nyange", []string{"Bweza", "Kabaya", "Nyange", "Rukore", "Rusasa"}},
		{"Musanze", "Remera", []string{"Gashaki", "Kabeza", "Kavumu", "Remera", "Rukore"}},
		{"Musanze", "Rwaza", []string{"Gakenke", "Kigombe", "Nyamigina", "Rwaza", "Uburondwe"}},
		{"Musanze", "Shingiro", []string{"Gashaki", "Kavumu", "Nyabageni", "Shingiro"}},

		// ════════════════════════════════════════
		// RULINDO DISTRICT
		// ════════════════════════════════════════
		{"Rulindo", "Base", []string{"Base", "Burega", "Cyungo", "Murambi"}},
		{"Rulindo", "Burega", []string{"Base", "Burega", "Murambi", "Ntarabana"}},
		{"Rulindo", "Bushoki", []string{"Bushoki", "Cyungo", "Murambi", "Shyorongi"}},
		{"Rulindo", "Buyoga", []string{"Buyoga", "Kisaro", "Masoro", "Tumba"}},
		{"Rulindo", "Cyinzuzi", []string{"Burega", "Cyinzuzi", "Murambi", "Ntarabana"}},
		{"Rulindo", "Cyungo", []string{"Base", "Bushoki", "Cyungo", "Murambi"}},
		{"Rulindo", "Kinihira", []string{"Cyungo", "Kinihira", "Murambi", "Shyorongi"}},
		{"Rulindo", "Kisaro", []string{"Buyoga", "Kisaro", "Masoro", "Ngoma"}},
		{"Rulindo", "Masoro", []string{"Buyoga", "Kisaro", "Masoro", "Tumba"}},
		{"Rulindo", "Mbogo", []string{"Mbogo", "Murambi", "Ntarabana", "Shyorongi"}},
		{"Rulindo", "Murambi", []string{"Base", "Burega", "Murambi", "Ntarabana"}},
		{"Rulindo", "Ngoma", []string{"Kisaro", "Masoro", "Ngoma", "Tumba"}},
		{"Rulindo", "Ntarabana", []string{"Burega", "Murambi", "Ntarabana", "Shyorongi"}},
		{"Rulindo", "Rukozo", []string{"Murambi", "Ntarabana", "Rukozo", "Shyorongi"}},
		{"Rulindo", "Rusiga", []string{"Kisaro", "Masoro", "Ngoma", "Rusiga"}},
		{"Rulindo", "Shyorongi", []string{"Bushoki", "Murambi", "Ntarabana", "Shyorongi"}},
		{"Rulindo", "Tumba", []string{"Buyoga", "Kisaro", "Masoro", "Tumba"}},

		// ════════════════════════════════════════
		// GISAGARA DISTRICT
		// ════════════════════════════════════════
		{"Gisagara", "Gikonko", []string{"Gikonko", "Kansi", "Muganza", "Ndora"}},
		{"Gisagara", "Gishubi", []string{"Gishubi", "Kigembe", "Muganza", "Ndora"}},
		{"Gisagara", "Kansi", []string{"Gikonko", "Kansi", "Muganza", "Save"}},
		{"Gisagara", "Kibirizi", []string{"Kibirizi", "Kigembe", "Muganza", "Musha"}},
		{"Gisagara", "Kigembe", []string{"Gishubi", "Kigembe", "Muganza", "Ndora"}},
		{"Gisagara", "Mamba", []string{"Gikonko", "Kansi", "Mamba", "Save"}},
		{"Gisagara", "Muganza", []string{"Gikonko", "Muganza", "Mugombwa", "Ndora"}},
		{"Gisagara", "Mugombwa", []string{"Muganza", "Mugombwa", "Musha", "Ndora"}},
		{"Gisagara", "Mukingo", []string{"Kigembe", "Mukingo", "Musha", "Ndora"}},
		{"Gisagara", "Musha", []string{"Kibirizi", "Muganza", "Musha", "Ndora"}},
		{"Gisagara", "Ndora", []string{"Gikonko", "Kigembe", "Muganza", "Ndora"}},
		{"Gisagara", "Nyanza", []string{"Gikonko", "Kansi", "Muganza", "Nyanza"}},
		{"Gisagara", "Save", []string{"Kansi", "Mamba", "Muganza", "Save"}},

		// ════════════════════════════════════════
		// HUYE DISTRICT
		// ════════════════════════════════════════
		{"Huye", "Gishamvu", []string{"Gishamvu", "Karama", "Mbazi", "Mukura"}},
		{"Huye", "Huye", []string{"Huye", "Karama", "Ngoma", "Tumba"}},
		{"Huye", "Karama", []string{"Gishamvu", "Karama", "Maraba", "Simbi"}},
		{"Huye", "Kigoma", []string{"Kigoma", "Maraba", "Ngoma", "Rwaniro"}},
		{"Huye", "Kinazi", []string{"Kinazi", "Mukura", "Rusatira", "Rwaniro"}},
		{"Huye", "Maraba", []string{"Karama", "Kigoma", "Maraba", "Simbi"}},
		{"Huye", "Mbazi", []string{"Gishamvu", "Mbazi", "Mukura", "Ruhashya"}},
		{"Huye", "Mukura", []string{"Kinazi", "Mbazi", "Mukura", "Rusatira"}},
		{"Huye", "Ngoma", []string{"Huye", "Kigoma", "Ngoma", "Tumba"}},
		{"Huye", "Ruhashya", []string{"Mbazi", "Mukura", "Ruhashya", "Rusatira"}},
		{"Huye", "Rusatira", []string{"Kinazi", "Mukura", "Ruhashya", "Rusatira"}},
		{"Huye", "Rwaniro", []string{"Kigoma", "Kinazi", "Rwaniro", "Simbi"}},
		{"Huye", "Simbi", []string{"Karama", "Maraba", "Rwaniro", "Simbi"}},
		{"Huye", "Tumba", []string{"Huye", "Kigoma", "Ngoma", "Tumba"}},

		// ════════════════════════════════════════
		// KAMONYI DISTRICT
		// ════════════════════════════════════════
		{"Kamonyi", "Gacurabwenge", []string{"Gacurabwenge", "Kayenzi", "Mugina", "Runda"}},
		{"Kamonyi", "Karama", []string{"Karama", "Kayenzi", "Musambira", "Rugarika"}},
		{"Kamonyi", "Kayenzi", []string{"Gacurabwenge", "Kayenzi", "Mugina", "Rukoma"}},
		{"Kamonyi", "Kayumbu", []string{"Kayumbu", "Mugina", "Musambira", "Runda"}},
		{"Kamonyi", "Mugina", []string{"Gacurabwenge", "Kayumbu", "Mugina", "Runda"}},
		{"Kamonyi", "Musambira", []string{"Karama", "Kayumbu", "Musambira", "Rugarika"}},
		{"Kamonyi", "Ngamba", []string{"Mugina", "Ngamba", "Nyarubaka", "Runda"}},
		{"Kamonyi", "Nyamiyaga", []string{"Kayenzi", "Mugina", "Nyamiyaga", "Rukoma"}},
		{"Kamonyi", "Nyarubaka", []string{"Mugina", "Ngamba", "Nyarubaka", "Rugarika"}},
		{"Kamonyi", "Rugarika", []string{"Karama", "Musambira", "Nyarubaka", "Rugarika"}},
		{"Kamonyi", "Rukoma", []string{"Kayenzi", "Nyamiyaga", "Rukoma", "Runda"}},
		{"Kamonyi", "Runda", []string{"Gacurabwenge", "Mugina", "Ngamba", "Runda"}},

		// ════════════════════════════════════════
		// MUHANGA DISTRICT
		// ════════════════════════════════════════
		{"Muhanga", "Cyeza", []string{"Cyeza", "Kabacuzi", "Kibangu", "Rongi"}},
		{"Muhanga", "Kabacuzi", []string{"Cyeza", "Kabacuzi", "Muhanga", "Shyogwe"}},
		{"Muhanga", "Kibangu", []string{"Cyeza", "Kibangu", "Muhanga", "Nyamabuye"}},
		{"Muhanga", "Kiyumba", []string{"Kabacuzi", "Kiyumba", "Muhanga", "Shyogwe"}},
		{"Muhanga", "Muhanga", []string{"Kibangu", "Muhanga", "Nyamabuye", "Nyarusange"}},
		{"Muhanga", "Mushishiro", []string{"Kabacuzi", "Muhanga", "Mushishiro", "Shyogwe"}},
		{"Muhanga", "Nyabinoni", []string{"Kabacuzi", "Nyabinoni", "Rongi", "Shyogwe"}},
		{"Muhanga", "Nyamabuye", []string{"Kibangu", "Muhanga", "Nyamabuye", "Nyarusange"}},
		{"Muhanga", "Nyarusange", []string{"Muhanga", "Nyamabuye", "Nyarusange", "Rongi"}},
		{"Muhanga", "Rongi", []string{"Cyeza", "Kabacuzi", "Nyabinoni", "Rongi"}},
		{"Muhanga", "Rugendabari", []string{"Kabacuzi", "Muhanga", "Rugendabari", "Shyogwe"}},
		{"Muhanga", "Shyogwe", []string{"Kabacuzi", "Kiyumba", "Mushishiro", "Shyogwe"}},

		// ════════════════════════════════════════
		// NYAMAGABE DISTRICT
		// ════════════════════════════════════════
		{"Nyamagabe", "Buruhukiro", []string{"Buruhukiro", "Gasaka", "Musebeya", "Tare"}},
		{"Nyamagabe", "Cyanika", []string{"Cyanika", "Gasaka", "Kibumbwe", "Musebeya"}},
		{"Nyamagabe", "Gasaka", []string{"Buruhukiro", "Gasaka", "Gatare", "Musebeya"}},
		{"Nyamagabe", "Gatare", []string{"Gasaka", "Gatare", "Kaduha", "Musange"}},
		{"Nyamagabe", "Kaduha", []string{"Gatare", "Kaduha", "Kamegeri", "Musange"}},
		{"Nyamagabe", "Kamegeri", []string{"Kaduha", "Kamegeri", "Kibirizi", "Musange"}},
		{"Nyamagabe", "Kibirizi", []string{"Kamegeri", "Kibirizi", "Mugano", "Musange"}},
		{"Nyamagabe", "Kibumbwe", []string{"Cyanika", "Gasaka", "Kibumbwe", "Kitabi"}},
		{"Nyamagabe", "Kitabi", []string{"Gasaka", "Kibumbwe", "Kitabi", "Musebeya"}},
		{"Nyamagabe", "Mbazi", []string{"Kaduha", "Mbazi", "Mugano", "Musange"}},
		{"Nyamagabe", "Mugano", []string{"Kibirizi", "Mbazi", "Mugano", "Musange"}},
		{"Nyamagabe", "Musange", []string{"Gatare", "Kaduha", "Kamegeri", "Musange"}},
		{"Nyamagabe", "Musebeya", []string{"Buruhukiro", "Gasaka", "Kitabi", "Musebeya"}},
		{"Nyamagabe", "Mushubi", []string{"Cyanika", "Gasaka", "Mushubi", "Tare"}},
		{"Nyamagabe", "Nkomane", []string{"Buruhukiro", "Musebeya", "Nkomane", "Tare"}},
		{"Nyamagabe", "Tare", []string{"Buruhukiro", "Musebeya", "Nkomane", "Tare"}},
		{"Nyamagabe", "Uwinkingi", []string{"Gasaka", "Kibumbwe", "Musebeya", "Uwinkingi"}},

		// ════════════════════════════════════════
		// NYANZA DISTRICT
		// ════════════════════════════════════════
		{"Nyanza", "Busasamana", []string{"Busasamana", "Kigoma", "Muyira", "Rwabicuma"}},
		{"Nyanza", "Kibirizi", []string{"Busoro", "Kibirizi", "Mukingo", "Ntyazo"}},
		{"Nyanza", "Busoro", []string{"Busoro", "Cyabakamyi", "Muyira", "Ntyazo"}},
		{"Nyanza", "Cyabakamyi", []string{"Busoro", "Cyabakamyi", "Kigoma", "Ntyazo"}},
		{"Nyanza", "Kigoma", []string{"Busasamana", "Cyabakamyi", "Kigoma", "Rwabicuma"}},
		{"Nyanza", "Mukingo", []string{"Kigoma", "Mukingo", "Muyira", "Nyagisozi"}},
		{"Nyanza", "Muyira", []string{"Busoro", "Mukingo", "Muyira", "Ntyazo"}},
		{"Nyanza", "Ntyazo", []string{"Busoro", "Cyabakamyi", "Muyira", "Ntyazo"}},
		{"Nyanza", "Nyagisozi", []string{"Kigoma", "Mukingo", "Nyagisozi", "Rwabicuma"}},
		{"Nyanza", "Rwabicuma", []string{"Busasamana", "Kigoma", "Nyagisozi", "Rwabicuma"}},

		// ════════════════════════════════════════
		// NYARUGURU DISTRICT
		// ════════════════════════════════════════
		{"Nyaruguru", "Busanze", []string{"Busanze", "Cyahinda", "Kibeho", "Munini"}},
		{"Nyaruguru", "Cyahinda", []string{"Busanze", "Cyahinda", "Mata", "Ngoma"}},
		{"Nyaruguru", "Kibeho", []string{"Busanze", "Kibeho", "Munini", "Ngera"}},
		{"Nyaruguru", "Kivu", []string{"Kivu", "Mata", "Munini", "Nyabimata"}},
		{"Nyaruguru", "Mata", []string{"Cyahinda", "Kivu", "Mata", "Ngoma"}},
		{"Nyaruguru", "Muganza", []string{"Muganza", "Munini", "Ngera", "Ruramba"}},
		{"Nyaruguru", "Munini", []string{"Busanze", "Kibeho", "Muganza", "Munini"}},
		{"Nyaruguru", "Ngera", []string{"Kibeho", "Muganza", "Munini", "Ngera"}},
		{"Nyaruguru", "Ngoma", []string{"Cyahinda", "Mata", "Ngoma", "Ruheru"}},
		{"Nyaruguru", "Nyabimata", []string{"Kivu", "Munini", "Nyabimata", "Nyagisozi"}},
		{"Nyaruguru", "Nyagisozi", []string{"Munini", "Nyabimata", "Nyagisozi", "Ruramba"}},
		{"Nyaruguru", "Ruheru", []string{"Mata", "Ngoma", "Ruheru", "Rusenge"}},
		{"Nyaruguru", "Ruramba", []string{"Muganza", "Munini", "Nyagisozi", "Ruramba"}},
		{"Nyaruguru", "Rusenge", []string{"Mata", "Ngoma", "Ruheru", "Rusenge"}},

		// ════════════════════════════════════════
		// RUHANGO DISTRICT
		// ════════════════════════════════════════
		{"Ruhango", "Bweramana", []string{"Bweramana", "Kabagali", "Kinazi", "Ruhango"}},
		{"Ruhango", "Byimana", []string{"Byimana", "Kabagali", "Kinihira", "Ntongwe"}},
		{"Ruhango", "Kabagali", []string{"Bweramana", "Byimana", "Kabagali", "Ruhango"}},
		{"Ruhango", "Kinazi", []string{"Bweramana", "Kinazi", "Mbuye", "Ruhango"}},
		{"Ruhango", "Kinihira", []string{"Byimana", "Kinihira", "Mwendo", "Ntongwe"}},
		{"Ruhango", "Mbuye", []string{"Bweramana", "Kinazi", "Mbuye", "Ruhango"}},
		{"Ruhango", "Mwendo", []string{"Kinihira", "Mbuye", "Mwendo", "Ntongwe"}},
		{"Ruhango", "Ntongwe", []string{"Byimana", "Kinihira", "Mwendo", "Ntongwe"}},
		{"Ruhango", "Ruhango", []string{"Bweramana", "Kabagali", "Kinazi", "Ruhango"}},

		// ════════════════════════════════════════
		// KARONGI DISTRICT
		// ════════════════════════════════════════
		{"Karongi", "Bwishyura", []string{"Bwishyura", "Gishari", "Gitesi", "Rubengera"}},
		{"Karongi", "Gishari", []string{"Bwishyura", "Gishari", "Mubuga", "Rubengera"}},
		{"Karongi", "Gishyita", []string{"Gishyita", "Mubuga", "Rugabano", "Ruganda"}},
		{"Karongi", "Gitesi", []string{"Bwishyura", "Gitesi", "Rubengera", "Twumba"}},
		{"Karongi", "Mubuga", []string{"Gishari", "Gishyita", "Mubuga", "Rugabano"}},
		{"Karongi", "Murambi", []string{"Murambi", "Murundi", "Rugabano", "Rwankuba"}},
		{"Karongi", "Murundi", []string{"Murambi", "Murundi", "Mutuntu", "Rwankuba"}},
		{"Karongi", "Mutuntu", []string{"Murundi", "Mutuntu", "Rugabano", "Twumba"}},
		{"Karongi", "Rubengera", []string{"Bwishyura", "Gitesi", "Rubengera", "Rugabano"}},
		{"Karongi", "Rugabano", []string{"Gishyita", "Mubuga", "Rugabano", "Ruganda"}},
		{"Karongi", "Ruganda", []string{"Gishyita", "Mubuga", "Rugabano", "Ruganda"}},
		{"Karongi", "Rwankuba", []string{"Murambi", "Murundi", "Rwankuba", "Twumba"}},
		{"Karongi", "Twumba", []string{"Gitesi", "Mutuntu", "Rwankuba", "Twumba"}},

		// ════════════════════════════════════════
		// NGORORERO DISTRICT
		// ════════════════════════════════════════
		{"Ngororero", "Bwira", []string{"Bwira", "Hindiro", "Muhanda", "Ngororero"}},
		{"Ngororero", "Gatumba", []string{"Gatumba", "Kabaya", "Matyazo", "Muhanda"}},
		{"Ngororero", "Hindiro", []string{"Bwira", "Hindiro", "Muhanda", "Ndaro"}},
		{"Ngororero", "Kabaya", []string{"Gatumba", "Kabaya", "Kageyo", "Matyazo"}},
		{"Ngororero", "Kageyo", []string{"Kabaya", "Kageyo", "Kavumu", "Matyazo"}},
		{"Ngororero", "Kavumu", []string{"Kageyo", "Kavumu", "Muhororo", "Nyange"}},
		{"Ngororero", "Matyazo", []string{"Gatumba", "Kabaya", "Kageyo", "Matyazo"}},
		{"Ngororero", "Muhanda", []string{"Bwira", "Gatumba", "Hindiro", "Muhanda"}},
		{"Ngororero", "Muhororo", []string{"Kavumu", "Muhororo", "Ndaro", "Ngororero"}},
		{"Ngororero", "Ndaro", []string{"Hindiro", "Muhororo", "Ndaro", "Ngororero"}},
		{"Ngororero", "Ngororero", []string{"Bwira", "Muhororo", "Ndaro", "Ngororero"}},
		{"Ngororero", "Nyange", []string{"Kavumu", "Muhororo", "Nyange", "Sovu"}},
		{"Ngororero", "Sovu", []string{"Kavumu", "Muhororo", "Nyange", "Sovu"}},

		// ════════════════════════════════════════
		// NYABIHU DISTRICT
		// ════════════════════════════════════════
		{"Nyabihu", "Bigogwe", []string{"Bigogwe", "Jenda", "Karago", "Mukamira"}},
		{"Nyabihu", "Jenda", []string{"Bigogwe", "Jenda", "Kabatwa", "Mukamira"}},
		{"Nyabihu", "Jomba", []string{"Jomba", "Karago", "Kintobo", "Mukamira"}},
		{"Nyabihu", "Kabatwa", []string{"Jenda", "Kabatwa", "Mukamira", "Muringa"}},
		{"Nyabihu", "Karago", []string{"Bigogwe", "Jomba", "Karago", "Mukamira"}},
		{"Nyabihu", "Kintobo", []string{"Jomba", "Karago", "Kintobo", "Rugera"}},
		{"Nyabihu", "Mukamira", []string{"Bigogwe", "Jenda", "Mukamira", "Rambura"}},
		{"Nyabihu", "Muringa", []string{"Kabatwa", "Mukamira", "Muringa", "Rambura"}},
		{"Nyabihu", "Rambura", []string{"Mukamira", "Muringa", "Rambura", "Shyira"}},
		{"Nyabihu", "Rugera", []string{"Karago", "Kintobo", "Rugera", "Rurembo"}},
		{"Nyabihu", "Rurembo", []string{"Karago", "Rugera", "Rurembo", "Shyira"}},
		{"Nyabihu", "Shyira", []string{"Rambura", "Rugera", "Rurembo", "Shyira"}},

		// ════════════════════════════════════════
		// NYAMASHEKE DISTRICT
		// ════════════════════════════════════════
		{"Nyamasheke", "Bushekeri", []string{"Bushekeri", "Cyato", "Kagano", "Shangi"}},
		{"Nyamasheke", "Bushenge", []string{"Bushenge", "Gihombo", "Macuba", "Mahembe"}},
		{"Nyamasheke", "Cyato", []string{"Bushekeri", "Cyato", "Kagano", "Rangiro"}},
		{"Nyamasheke", "Gihombo", []string{"Bushenge", "Gihombo", "Kanjongo", "Macuba"}},
		{"Nyamasheke", "Kagano", []string{"Bushekeri", "Cyato", "Kagano", "Karambi"}},
		{"Nyamasheke", "Kanjongo", []string{"Gihombo", "Kanjongo", "Karengera", "Macuba"}},
		{"Nyamasheke", "Karambi", []string{"Kagano", "Karambi", "Kirimbi", "Rangiro"}},
		{"Nyamasheke", "Karengera", []string{"Kanjongo", "Karengera", "Kirimbi", "Nyabitekeri"}},
		{"Nyamasheke", "Kirimbi", []string{"Karambi", "Karengera", "Kirimbi", "Rangiro"}},
		{"Nyamasheke", "Macuba", []string{"Bushenge", "Gihombo", "Kanjongo", "Macuba"}},
		{"Nyamasheke", "Mahembe", []string{"Bushenge", "Macuba", "Mahembe", "Shangi"}},
		{"Nyamasheke", "Nyabitekeri", []string{"Karengera", "Nyabitekeri", "Rangiro", "Ruharambuga"}},
		{"Nyamasheke", "Rangiro", []string{"Cyato", "Karambi", "Kirimbi", "Rangiro"}},
		{"Nyamasheke", "Ruharambuga", []string{"Nyabitekeri", "Rangiro", "Ruharambuga", "Shangi"}},
		{"Nyamasheke", "Shangi", []string{"Bushekeri", "Mahembe", "Ruharambuga", "Shangi"}},

		// ════════════════════════════════════════
		// RUBAVU DISTRICT
		// ════════════════════════════════════════
		{"Rubavu", "Bugeshi", []string{"Bugeshi", "Busasamana", "Cyanzarwe", "Mudende"}},
		{"Rubavu", "Busasamana", []string{"Bugeshi", "Busasamana", "Cyanzarwe", "Kanama"}},
		{"Rubavu", "Cyanzarwe", []string{"Bugeshi", "Busasamana", "Cyanzarwe", "Mudende"}},
		{"Rubavu", "Gisenyi", []string{"Gisenyi", "Kanzenze", "Rubavu", "Rugerero"}},
		{"Rubavu", "Kanama", []string{"Busasamana", "Kanama", "Mudende", "Nyamyumba"}},
		{"Rubavu", "Kanzenze", []string{"Gisenyi", "Kanzenze", "Nyakiriba", "Rubavu"}},
		{"Rubavu", "Mudende", []string{"Bugeshi", "Cyanzarwe", "Kanama", "Mudende"}},
		{"Rubavu", "Nyakiriba", []string{"Kanzenze", "Nyakiriba", "Nyamyumba", "Rubavu"}},
		{"Rubavu", "Nyamyumba", []string{"Kanama", "Nyakiriba", "Nyamyumba", "Nyundo"}},
		{"Rubavu", "Nyundo", []string{"Nyamyumba", "Nyundo", "Rubavu", "Rugerero"}},
		{"Rubavu", "Rubavu", []string{"Gisenyi", "Kanzenze", "Rubavu", "Rugerero"}},
		{"Rubavu", "Rugerero", []string{"Gisenyi", "Nyundo", "Rubavu", "Rugerero"}},

		// ════════════════════════════════════════
		// RUSIZI DISTRICT
		// ════════════════════════════════════════
		{"Rusizi", "Bugarama", []string{"Bugarama", "Gashonga", "Giheke", "Nkungu"}},
		{"Rusizi", "Butare", []string{"Butare", "Giheke", "Gihundwe", "Nkungu"}},
		{"Rusizi", "Bweyeye", []string{"Bweyeye", "Gashonga", "Gikundamvura", "Nzahaha"}},
		{"Rusizi", "Gashonga", []string{"Bugarama", "Gashonga", "Mururu", "Nzahaha"}},
		{"Rusizi", "Giheke", []string{"Bugarama", "Butare", "Giheke", "Nkungu"}},
		{"Rusizi", "Gihundwe", []string{"Butare", "Gihundwe", "Kamembe", "Nyakabuye"}},
		{"Rusizi", "Gikundamvura", []string{"Bweyeye", "Gikundamvura", "Nkanka", "Nzahaha"}},
		{"Rusizi", "Gitambi", []string{"Bugarama", "Gashonga", "Gitambi", "Muganza"}},
		{"Rusizi", "Kamembe", []string{"Gihundwe", "Kamembe", "Nyakabuye", "Nyakarenzo"}},
		{"Rusizi", "Muganza", []string{"Bugarama", "Gashonga", "Muganza", "Nkungu"}},
		{"Rusizi", "Mururu", []string{"Gashonga", "Muganza", "Mururu", "Nzahaha"}},
		{"Rusizi", "Nkanka", []string{"Gikundamvura", "Nkanka", "Nkombo", "Nzahaha"}},
		{"Rusizi", "Nkombo", []string{"Nkanka", "Nkombo", "Nkungu", "Nzahaha"}},
		{"Rusizi", "Nkungu", []string{"Bugarama", "Giheke", "Muganza", "Nkungu"}},
		{"Rusizi", "Nyakabuye", []string{"Gihundwe", "Kamembe", "Nyakabuye", "Nyakarenzo"}},
		{"Rusizi", "Nyakarenzo", []string{"Kamembe", "Nyakabuye", "Nyakarenzo", "Rwimbogo"}},
		{"Rusizi", "Nzahaha", []string{"Bweyeye", "Gashonga", "Mururu", "Nzahaha"}},
		{"Rusizi", "Rwimbogo", []string{"Kamembe", "Nyakarenzo", "Nzahaha", "Rwimbogo"}},

		// ════════════════════════════════════════
		// RUTSIRO DISTRICT
		// ════════════════════════════════════════
		{"Rutsiro", "Boneza", []string{"Boneza", "Gihango", "Musasa", "Rusebeya"}},
		{"Rutsiro", "Gihango", []string{"Boneza", "Gihango", "Kivumu", "Mushonyi"}},
		{"Rutsiro", "Kigeyo", []string{"Kigeyo", "Manihira", "Mushonyi", "Nyabirasi"}},
		{"Rutsiro", "Kivumu", []string{"Gihango", "Kivumu", "Murunda", "Mushonyi"}},
		{"Rutsiro", "Manihira", []string{"Kigeyo", "Manihira", "Mukura", "Nyabirasi"}},
		{"Rutsiro", "Mukura", []string{"Manihira", "Mukura", "Mushubati", "Nyabirasi"}},
		{"Rutsiro", "Murunda", []string{"Kivumu", "Murunda", "Musasa", "Mushonyi"}},
		{"Rutsiro", "Musasa", []string{"Boneza", "Murunda", "Musasa", "Rusebeya"}},
		{"Rutsiro", "Mushonyi", []string{"Gihango", "Kigeyo", "Kivumu", "Mushonyi"}},
		{"Rutsiro", "Mushubati", []string{"Mukura", "Mushubati", "Nyabirasi", "Ruhango"}},
		{"Rutsiro", "Nyabirasi", []string{"Kigeyo", "Manihira", "Mukura", "Nyabirasi"}},
		{"Rutsiro", "Ruhango", []string{"Mushubati", "Nyabirasi", "Ruhango", "Rusebeya"}},
		{"Rutsiro", "Rusebeya", []string{"Boneza", "Musasa", "Ruhango", "Rusebeya"}},
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
		{"Gasabo", "Kimironko", "Kimironko", []string{"Ituze", "Amahoro", "Umwezi", "Imboni"}},
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
