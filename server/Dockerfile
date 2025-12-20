# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy any additional necessary files (if needed)
# COPY internal/DB/schema.sql /app/schema.sql
# COPY internal/DB/query.sql /app/query.sql

# Expose ports (will be overridden by docker-compose)
EXPOSE 3001 3002 3003

# CMD will be specified in docker-compose.yml 