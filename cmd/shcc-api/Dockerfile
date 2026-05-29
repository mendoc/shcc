# Utiliser une version récente de Go
FROM golang:1.23-alpine AS builder

WORKDIR /src

# Copier les fichiers go.mod et go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copier tout le monorepo
COPY . .

# Compiler l'API spécifiquement
RUN CGO_ENABLED=0 GOOS=linux go build -o /shcc-api ./cmd/shcc-api/main.go

# Image finale légère
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /shcc-api .
COPY --from=builder /src/cmd/shcc-api/install.sh.tmpl ./cmd/shcc-api/install.sh.tmpl
COPY --from=builder /src/bin ./bin

# Variable d'environnement par défaut
ENV PORT=8080

EXPOSE 8080

CMD ["./shcc-api"]
