# ---------- Build stage ----------
# AS builder “multi-stage build” : on va construire le binaire, puis l’utiliser dans une autre image plus légère.
FROM golang:1.24.2-alpine AS builder 

# Dépendances système pour compiler Go et SQLite
# apk gestionnaire de package alpine -> no cache pour ne pas garder les fichier de telechargement des dependences (image plus legere)
# Ces paquets sont nécessaires uniquement pendant la compilation (builder stage)
# Une fois le binaire Go compilé, tu n’as plus besoin de gcc ni de musl-dev pour exécuter ton application
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Définit le dossier de travail à l’intérieur du conteneur.
# Tous les fichiers copiés ou commandes suivantes vont se faire dans /app.
WORKDIR /app

# Copier et télécharger les dépendances Go
COPY go.mod go.sum ./
RUN go mod download

# Copie tout le projet dans le conteneur, dans /app
COPY . .

# run lance le compile de l'application, CGO dit a go d'inclure la partie C pour sqlite. goos compile pour linux et nom du binaire final /forum
RUN CGO_ENABLED=1 GOOS=linux go build -o /forum

# On cree une image tres legere basé seulement sur alpine linux (pour ne pas inclure go et tout les outils de compilation dans limage final)
FROM alpine:3.20

# Installe SQLite dans l’image finale
RUN apk add --no-cache sqlite

# Copier le binaire compilé
# 1er /forum (source) → c’est le binaire Go compilé dans le conteneur builder.
# 2eme /forum (destination) → emplacement dans l’image finale.
COPY --from=builder /forum /forum

# Copier templates et static
# Ces fichiers ne sont pas compilés, mais l'application web en a besoin à runtime (utiliser par le programme lorsqu'il tourne).
COPY --from=builder /app/templates /templates
COPY --from=builder /app/static /static

# Créer dossier pour la DB et config utilisateur non-root
# Crée un dossier pour les données persistantes
RUN mkdir -p /data \
# Crée un utilisateur non-root pour plus de sécurité
    && adduser -D appuser \
#  Donne la propriété du dossier /data à cet utilisateur
    && chown -R appuser /data
# Informer le programme que toute commandes seront executer par appuser et pas root (appuser :utilisateur d'application)
USER appuser

# Définir le dossier DB persistant
VOLUME /data

# Exposer le port
# Informe Docker que le conteneur écoute sur le port 5090.
EXPOSE 5090

# Lancer le binaire
# Définit le programme qui se lance automatiquement quand le conteneur démarre
ENTRYPOINT ["/forum"]
