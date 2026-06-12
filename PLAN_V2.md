# rwandapi v2 Plan

**Theme: Live data, authentication, and deeper coverage**

## New Data Endpoints

1. **Transportation** — Bus routes (Kigali Bus Services, Royal Express, Volcano Express), taxi parks, moto fares by zone, KBS route maps
2. **Education** — Universities, TVET schools, primary/secondary schools by district, with enrollment stats
3. **Healthcare** — Hospitals, health centers, pharmacies by district, referral hospital contacts
4. **Tourism** — National parks (Volcanoes, Akagera, Nyungwe), museums, genocide memorials, entry fees
5. **Media** — Radio stations, TV channels, newspapers with frequencies/URLs
6. **Postal Codes** — Rwanda postal code lookup by district/sector

## Technical Improvements

7. **Live exchange rates** — Scrape BNR daily rates instead of static seed data
8. **Pagination** — `?page=1&limit=50` on all list endpoints
9. **Rate limiting** — Per-IP throttling to prevent abuse
10. **API versioning** — Proper `/v2` prefix with breaking changes, `/v1` stays stable
11. **OpenAPI/Swagger docs** — Auto-generated API documentation at `/docs`
12. **Health check** — `/health` endpoint for uptime monitoring
13. **CORS configuration** — Configurable allowed origins
14. **Caching headers** — ETags and Cache-Control for static data
15. **Docker Compose** — Dev setup with hot reload

## Stretch Goals

16. **Authentication** — Optional API keys for higher rate limits
17. **Admin endpoints** — CRUD for managing seed data without redeploying
18. **Webhook subscriptions** — Notify when exchange rates update
19. **Full village data** — Complete 14,837 villages (currently only sample for Kimironko)
20. **Complete cell data** — Fill remaining ~426 cells to reach the full 2,148

## Priority

The most impactful v2 items are **pagination + rate limiting + live exchange rates + Swagger docs** — they make the API production-ready. Transportation and tourism data would be the most interesting new endpoints.
