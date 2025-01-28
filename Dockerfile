# Builder stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for better cache
COPY go.mod go.sum ./

# Install dependencies
RUN go mod tidy

# Copy the rest of the source code
COPY . ./

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Build the main application binary
RUN go build -o main cmd/main/main.go

# Final stage (runtime)
FROM alpine:latest

# Install necessary certificates and dependencies
RUN apk --no-cache add ca-certificates

# Copy the built binary and goose from the builder image
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/main /main
COPY --from=builder /app/migrations /migrations

# Make the binary executable
RUN chmod +x /main

# Define the entry point
CMD ["/main"]
