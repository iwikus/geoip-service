# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o geoip-service .

# Stage 2: Runtime with embedded databases
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/geoip-service .

# Embed all three databases (downloaded during CI before docker build)
RUN mkdir -p /data
COPY data/GeoLite2-City.mmdb    /data/GeoLite2-City.mmdb
COPY data/GeoLite2-Country.mmdb /data/GeoLite2-Country.mmdb
COPY data/GeoLite2-ASN.mmdb     /data/GeoLite2-ASN.mmdb

# Default symlink so the default CMD works without arguments
RUN ln -s /data/GeoLite2-City.mmdb /data/geodb.mmdb

EXPOSE 5000

# Default uses City DB. Override at runtime:
#   docker run ... -e DB=/data/GeoLite2-Country.mmdb
#   docker run ... /app/geoip-service -db=/data/GeoLite2-ASN.mmdb -lookup=asn
CMD ["/app/geoip-service", "-db=/data/geodb.mmdb"]
