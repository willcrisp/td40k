# Build stage
FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copy all source first so go mod tidy can resolve everything
COPY . .

# Generate a valid go.sum and download deps
RUN GONOSUMDB=* GOFLAGS=-mod=mod go mod tidy

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server .

# Run stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
