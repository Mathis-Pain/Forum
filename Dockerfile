# ---------- Build stage ----------
FROM golang:1.24.2-alpine AS builder

# Dépendances système pour compiler Go et SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copier et télécharger les dépendances Go
COPY go.mod go.sum ./
RUN go mod download

# Copier le reste du projet
COPY . .

# Compiler le binaire Linux statique
RUN CGO_ENABLED=1 GOOS=linux go build -o /forum

# ---------- Final stage ----------
FROM alpine:3.20

# Installer SQLite (utile si ton code l'utilise en runtime)
RUN apk add --no-cache sqlite

# Copier le binaire compilé
COPY --from=builder /forum /forum

# Copier templates et static
COPY --from=builder /app/templates /templates
COPY --from=builder /app/static /static

# Créer dossier pour la DB et config utilisateur non-root
RUN mkdir -p /data \
    && adduser -D appuser \
    && chown -R appuser /data

USER appuser

# Définir le dossier DB persistant
VOLUME /data

# Exposer le port
EXPOSE 5090

# Lancer le binaire
ENTRYPOINT ["/forum"]
