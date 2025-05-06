# Lead Locate

Google API documentation: [Google Maps API](https://serpapi.com/google-maps-api).
SerpApi documentation for Go: [SerpApi - Golang](https://serpapi.com/integrations/go).

```bash
curl --get https://serpapi.com/search \
    -d engine="google_maps" \
    -d q="Coffee" \
    -d ll="@40.7455096,-74.0083012,14z" \
    -d api_key="<API KEY>"
```

Nominatim API: [Nominatim 5.0.0 Manua](https://nominatim.org/release-docs/latest/api/Search/#structured-query)

```bash
curl "https://nominatim.openstreetmap.org/search?q=Treviso,%20Veneto&format=json&limit=3" | jq
```

## Docker

Build the environment:

```bash
docker compose -f deploy/docker-compose.yaml build
```

Then run it:

```bash
docker compose -f deploy/docker-compose.yaml up -d
```

---
