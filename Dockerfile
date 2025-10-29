FROM golang:1.24.2-alpine

WORKDIR /app

# Installer les dépendances nécessaires pour CGO et SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copier tout le projet
COPY . .

# Télécharger les dépendances
RUN go mod download

# Build avec CGO activé
RUN CGO_ENABLED=1 go build -o server .

EXPOSE 5080

CMD ["./server"]