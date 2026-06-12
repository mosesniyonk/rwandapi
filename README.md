# rwandapi

A free, open-source REST API serving Rwanda data — administrative divisions, banks, telecoms, mobile money, emergency numbers, holidays, and more.

No API key required. No rate-limit gotchas. Just clean JSON.

## Base URL

```
https://rwandapi.onrender.com
```

For local development: `http://localhost:8080`

## Endpoints

### Geography

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/provinces` | List all 5 provinces |
| `GET` | `/api/v1/provinces/:id` | Single province with its districts |
| `GET` | `/api/v1/districts` | List all 30 districts (filter: `?province_id=`) |
| `GET` | `/api/v1/districts/:id` | Single district with its sectors |
| `GET` | `/api/v1/sectors` | List all 416 sectors (filter: `?district_id=`) |
| `GET` | `/api/v1/sectors/:id` | Single sector with its cells |
| `GET` | `/api/v1/cells` | List all 1,722 cells (filter: `?sector_id=`) |
| `GET` | `/api/v1/cells/:id` | Single cell with its villages |
| `GET` | `/api/v1/villages` | List villages (filter: `?cell_id=`) |
| `GET` | `/api/v1/search?q=kimironko` | Search across all geographic entities |

### Finance

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/banks` | List 10 banks with 129 branches, SWIFT codes, contacts |
| `GET` | `/api/v1/banks/:id` | Single bank with all its branches |
| `GET` | `/api/v1/exchange-rates` | BNR exchange rates for 18 currencies against RWF |
| `GET` | `/api/v1/mobile-money` | Mobile money services (MoMo, Airtel Money) with fee schedules |
| `GET` | `/api/v1/mobile-money/:id` | Single service with fees and limits |

### Telecoms & ISPs

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/telecoms` | List all telecom providers and ISPs (filter: `?type=mobile` or `?type=isp`) |
| `GET` | `/api/v1/telecoms/:id` | Single provider with plans and USSD codes |

### Utilities

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/utilities/electricity/tariffs` | EUCL electricity tariff bands |
| `GET` | `/api/v1/utilities/water/tariffs` | WASAC water tariff bands |

### General

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/info` | Country metadata (capital, population, languages, etc.) |
| `GET` | `/api/v1/holidays` | 14 national public holidays |
| `GET` | `/api/v1/emergency-numbers` | 20 emergency and useful numbers (filter: `?category=emergency`) |

## Example Responses

### GET /api/v1/provinces

```json
{
  "success": true,
  "data": [
    { "id": "a1b2c3...", "name": "Eastern Province", "created_at": "2024-01-01T00:00:00Z" },
    { "id": "d4e5f6...", "name": "Kigali City", "created_at": "2024-01-01T00:00:00Z" }
  ],
  "count": 5
}
```

### GET /api/v1/telecoms?type=mobile

```json
{
  "success": true,
  "data": [
    {
      "id": "f7g8h9...",
      "name": "MTN Rwanda",
      "short_name": "MTN",
      "type": "mobile",
      "website": "https://www.mtn.co.rw",
      "customer_care": "456",
      "ussd_codes": [
        { "code": "*182#", "description": "Main Menu" },
        { "code": "*131#", "description": "Check Balance" }
      ],
      "plans": [
        { "id": "...", "name": "Daily 100MB", "type": "data", "data_cap": "100MB", "price": "200", "validity": "1 day" }
      ]
    }
  ],
  "count": 2
}
```

### GET /api/v1/emergency-numbers?category=emergency

```json
{
  "success": true,
  "data": [
    { "id": "...", "name": "Rwanda National Police", "number": "112", "category": "emergency", "description": "Police emergency line" },
    { "id": "...", "name": "Ambulance / SAMU", "number": "912", "category": "emergency", "description": "Medical emergency and ambulance service" }
  ],
  "count": 7
}
```

### GET /api/v1/search?q=kimironko

```json
{
  "success": true,
  "data": [
    { "type": "sector", "id": "abc...", "name": "Kimironko", "path": "Kigali City > Gasabo > Kimironko" },
    { "type": "cell", "id": "def...", "name": "Kimironko", "path": "Kigali City > Gasabo > Kimironko > Kimironko" }
  ],
  "count": 2
}
```

## Data Coverage

| Entity | Count |
|--------|-------|
| Provinces | 5 |
| Districts | 30 |
| Sectors | 416 |
| Cells | 1,722 |
| Banks | 10 (129 branches) |
| Exchange Rates | 18 currencies |
| Telecom Providers | 5 (MTN, Airtel, BSC, Liquid, Canal Box) |
| Mobile Money | 2 (MoMo, Airtel Money) |
| Emergency Numbers | 20 |
| Public Holidays | 14 |

## Self-Hosting

### From source

```bash
git clone https://github.com/mosesniyonk/rwandapi.git
cd rwandapi
go build -o rwandapi .
./rwandapi
```

The server starts on port 8080 by default. Set the `PORT` environment variable to change it.

### Docker

```bash
docker build -t rwandapi .
docker run -p 8080:8080 rwandapi
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `RWANDAPI_DB_PATH` | `rwandapi.db` | Path to SQLite database file |

## Data Sources

- **Administrative divisions**: National Institute of Statistics of Rwanda (NISR)
- **Bank data**: National Bank of Rwanda (BNR)
- **Electricity tariffs**: Rwanda Utilities Regulatory Authority (RURA) / EUCL
- **Water tariffs**: Rwanda Utilities Regulatory Authority (RURA) / WASAC
- **Telecom plans**: MTN Rwanda, Airtel Rwanda, BSC, Liquid, Canal Box
- **Mobile money fees**: MTN MoMo, Airtel Money official fee schedules

> Data is seeded at startup. Rates and fees may change over time. Contributions to keep data current are welcome.

## Tech Stack

- **Go** with [Chi](https://github.com/go-chi/chi) router
- **SQLite** (embedded, zero external dependencies)
- Fully self-contained single binary

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/add-tourism-data`)
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

See [PLAN_V2.md](PLAN_V2.md) for the v2 roadmap and areas where contributions are welcome.

## License

MIT License - see [LICENSE](LICENSE) for details.
