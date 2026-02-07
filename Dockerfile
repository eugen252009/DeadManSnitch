# STAGE 1: Bauen des Binaries
FROM golang:1.25-alpine AS builder

# Arbeitsverzeichnis im Container
WORKDIR /app

# Abhängigkeiten kopieren und installieren
COPY go.mod ./
# Falls du eine go.sum hast, aktiviere die nächste Zeile:
# COPY go.sum ./
RUN go mod download

# Quellcode kopieren
COPY . .

# Programm statisch kompilieren für maximale Kompatibilität
RUN go build -o snitch-monitor .

# STAGE 2: Das schlanke Laufzeit-Image
FROM alpine:latest

# WICHTIG: Zertifikate installieren, damit HTTPS (Discord & Webseiten) funktioniert
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Nur das fertige Programm vom Builder-Image kopieren
COPY --from=builder /app/snitch-monitor .

# Programm starten
CMD ["./snitch-monitor"]