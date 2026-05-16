# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o pqr-server ./cmd/pqr/main.go

# Final Stage
FROM alpine:latest

WORKDIR /app
RUN apk --no-cache add ca-certificates curl
COPY --from=builder /app/pqr-server .
COPY --from=builder /app/web ./web
COPY --from=builder /app/docs ./docs

# Expose the API port
EXPOSE 8196

# Environment variables
ENV PORT=8196
ENV DATABASE_URL=postgresql://root@localhost:26257/antigravity?sslmode=disable

# Run the server
CMD ["./pqr-server"]
