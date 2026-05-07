# Stage 1: Build
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache curl ca-certificates

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o geoip-service .

# Download databases

RUN curl -fsSL \
      "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb" \
      -o GeoLite2-City.mmdb && \
    curl -fsSL \
      "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-Country.mmdb" \
      -o GeoLite2-Country.mmdb && \
    curl -fsSL \
      "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-ASN.mmdb" \
      -o GeoLite2-ASN.mmdb

# Stage 2: Runtime with embedded databases
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/geoip-service .
COPY --from=builder /app/*.mmdb .

EXPOSE 5000

CMD ["/app/geoip-service" ]

