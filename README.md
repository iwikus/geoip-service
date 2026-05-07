# geoip-service (maintained fork)

Fork of [klauspost/geoip-service](https://github.com/klauspost/geoip-service) — the original was archived in 2019.

**Changes in this fork:**
- Modern multi-stage `Dockerfile` (Go 1.22 + Alpine 3.19, ~15 MB image)
- GitHub Actions workflow that auto-builds weekly with fresh [GeoLite2 databases](https://github.com/P3TERX/GeoLite.mmdb)
- Image published to GitHub Container Registry (`ghcr.io`)

## Quick start

```bash
# Pull latest (databases embedded)
docker pull ghcr.io/iwikus/geoip-service:latest

# Run on port 5000
docker run --rm -p 5000:5000 ghcr.io/iwikus/geoip-service:latest
```

Query: `http://localhost:5000/1.2.3.4`

## Mount your own database instead

```bash
docker run --rm \
  -p 5000:5000 \
  -v /path/to/GeoLite2-City.mmdb:/data/geodb.mmdb \
  ghcr.io/iwikus/geoip-service:latest
```

## Available databases in the image

| Path | Type |
|------|------|
| `/data/GeoLite2-City.mmdb` | City (default) |
| `/data/GeoLite2-Country.mmdb` | Country |
| `/data/GeoLite2-ASN.mmdb` | ASN |

Switch database at runtime:
```bash
docker run --rm -p 5000:5000 ghcr.io/iwikus/geoip-service:latest \
  /app/geoip-service -db=/data/GeoLite2-Country.mmdb -lookup=country
```

## Auto-update schedule

The GitHub Actions workflow rebuilds the image every **Monday at 03:00 UTC**, picking up the latest DB snapshot from [P3TERX/GeoLite.mmdb](https://github.com/P3TERX/GeoLite.mmdb). You can also trigger it manually from the Actions tab.

## Service options

```
-db string       Path to .mmdb file (default "GeoLite2-City.mmdb")
-listen string   Listen address (default ":5000")
-lookup string   "city" or "country" (default "city")
-pretty bool     Pretty-print JSON output
-threads int     Worker threads (default: CPU count)
-cache int       Cache TTL in seconds, 0 = disabled
```
