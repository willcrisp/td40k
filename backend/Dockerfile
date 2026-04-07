# Build stage
FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copy go module files first for layer caching
COPY backend/go.mod backend/go.sum ./
RUN GONOSUMDB=* GOFLAGS=-mod=mod go mod download

# Copy backend source and migration files
COPY backend/ ./
COPY db/migrations/ ./db/migrations/

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# Run stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/db/migrations/ ./db/migrations/

EXPOSE 8080

CMD ["./server"]
