# Build Stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
# Install build dependencies for SQLite (CGO)
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Build static binary
RUN CGO_ENABLED=1 GOOS=linux go build -o server cmd/server/main.go

# Runtime Stage
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .

# Expose port (Render sets PORT env var, but good to document)
EXPOSE 8002

CMD ["./server"]
