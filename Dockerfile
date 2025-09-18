# Build stage
FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /forum

# Final stage
FROM alpine:3.20

# Copier le binaire
COPY --from=builder /forum /forum

# Dossier pour la DB
RUN mkdir -p /data \
    && adduser -D appuser \
    && chown -R appuser /data

USER appuser

VOLUME /data
EXPOSE 5080
ENTRYPOINT ["/forum"]
