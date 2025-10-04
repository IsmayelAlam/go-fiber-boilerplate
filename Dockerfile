FROM golang:1.25.1-alpine3.21 AS builder

WORKDIR /app

# Install packages 
RUN go install github.com/air-verse/air@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Copy go mod files FIRST
COPY go.mod go.sum ./
RUN go mod download

# Ensure build dir exists
RUN mkdir -p build

EXPOSE 8080

# Run air
CMD ["air"]