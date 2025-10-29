FROM golang:1.24.2-alpine

WORKDIR /app

# Installer les dépendances nécessaires pour CGO et SQLite
# apk alpinep package keeper gestionaire de pacquet  (--no-cache : pour ne pas garder un cache des fichier telechargé )
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copier tout le projet
COPY . .

# Télécharger les dépendances
RUN go mod download

# Build avec CGO activé
# Outil de Go qui permet d’utiliser du code C dans Go (CGO_ENABLED est une variable d’environnement qui contrôle si l’on utilise CGO. util pour SQLITE)
RUN CGO_ENABLED=1 go build -o server ./cmd/server

EXPOSE 5080

CMD ["./server"]