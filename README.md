# rwandapi 🇷🇼

A free, open-source REST API for Rwanda administrative divisions, banks, and utility data.

No API key required. No rate-limit gotchas. Just clean JSON.

## Base URL

```
https://rwandapi.onrender.com
```

For local development: `http://localhost:8080`

## Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/provinces` | List all 5 provinces |
| `GET` | `/api/v1/provinces/:id` | Single province with its districts |
| `GET` | `/api/v1/districts` | List all 30 districts (filter: `?province_id=1`) |
| `GET` | `/api/v1/districts/:id` | Single district with its sectors |
| `GET` | `/api/v1/sectors` | List sectors (filter: `?district_id=1`) |
| `GET` | `/api/v1/sectors/:id` | Single sector with its cells |
| `GET` | `/api/v1/cells` | List cells (filter: `?sector_id=1`) |
| `GET` | `/api/v1/cells/:id` | Single cell with its villages |
| `GET` | `/api/v1/villages` | List villages (filter: `?cell_id=1`) |
| `GET` | `/api/v1/banks` | List Rwanda banks with SWIFT codes and branches |
| `GET` | `/api/v1/utilities/electricity/tariffs` | EUCL electricity tariff bands |
| `GET` | `/api/v1/utilities/water/tariffs` | WASAC water tariff bands |
| `GET` | `/api/v1/search?q=kimironko` | Search across all geographic entities |
| `GET` | `/health` | Health check |

## Example Responses

### GET /api/v1/provinces

```json
{
  "success": true,
  "data": [
    { "id": 1, "name": "Eastern Province", "created_at": "2024-01-01T00:00:00Z" },
    { "id": 2, "name": "Kigali City", "created_at": "2024-01-01T00:00:00Z" },
    { "id": 3, "name": "Northern Province", "created_at": "2024-01-01T00:00:00Z" },
    { "id": 4, "name": "Southern Province", "created_at": "2024-01-01T00:00:00Z" },
    { "id": 5, "name": "Western Province", "created_at": "2024-01-01T00:00:00Z" }
  ],
  "count": 5
}
```

### GET /api/v1/districts?province_id=2

```json
{
  "success": true,
  "data": [
    { "id": 1, "province_id": 2, "name": "Gasabo", "created_at": "2024-01-01T00:00:00Z" },
    { "id": 2, "province_id": 2, "name": "Kicukiro", "created_at": "2024-01-01T00:00:00Z" },
    { "id": 3, "province_id": 2, "name": "Nyarugenge", "created_at": "2024-01-01T00:00:00Z" }
  ],
  "count": 3
}
```

### GET /api/v1/search?q=kimironko

```json
{
  "success": true,
  "data": [
    {
      "type": "sector",
      "id": 9,
      "name": "Kimironko",
      "path": "Kigali City > Gasabo > Kimironko"
    },
    {
      "type": "cell",
      "id": 3,
      "name": "Kimironko",
      "path": "Kigali City > Gasabo > Kimironko > Kimironko"
    }
  ],
  "count": 2
}
```

## Rate Limiting

- **100 requests per minute** per IP address
- When exceeded, the API returns `429 Too Many Requests` with a `Retry-After: 60` header

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

> Note: Data is seeded at startup. Tariff rates and bank information may change over time. Contributions to keep data current are welcome.

## Tech Stack

- **Go** with [Chi](https://github.com/go-chi/chi) router
- **SQLite** (embedded, zero external dependencies)
- Fully self-contained single binary

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/add-more-sectors`)
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

Areas where contributions are especially welcome:
- Adding more sectors, cells, and villages
- Updating tariff data
- Adding new data endpoints (exchange rates, holidays, etc.)

## License

MIT License - see [LICENSE](LICENSE) for details.
