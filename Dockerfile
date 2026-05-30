# Utiliser une version récente de Go (1.24+ requis par certaines dépendances)
FROM golang:1.24-alpine AS builder

WORKDIR /src

# Copier les fichiers go.mod et go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copier uniquement les dossiers nécessaires à l'API (excluant cmd/shcc-cli/)
COPY cmd/shcc-api/ ./cmd/shcc-api/
COPY internal/ ./internal/
COPY install.sh.tmpl ./
COPY bin/ ./bin/

# Compiler l'API spécifiquement
RUN CGO_ENABLED=0 GOOS=linux go build -o /shcc-api ./cmd/shcc-api/main.go

# Image finale légère
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /shcc-api .
COPY install.sh.tmpl ./install.sh.tmpl
COPY bin/ ./bin/
COPY migrations/ ./migrations/

# Variable d'environnement par défaut
ENV PORT=8080

EXPOSE 8080

CMD ["./shcc-api"]
