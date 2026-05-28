# Utiliser Go 1.22 pour la compilation
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copier les fichiers go.mod et go.sum
COPY go.mod ./
# RUN go mod download (Activer si vous avez des dépendances externes)

# Copier tout le monorepo
COPY . .

# Compiler l'API spécifiquement
RUN CGO_ENABLED=0 GOOS=linux go build -o /shcc-api ./cmd/shcc-api/main.go

# Image finale légère
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /shcc-api /shcc-api

# Variable d'environnement par défaut
ENV PORT 8080

EXPOSE 8080

CMD ["/shcc-api"]
