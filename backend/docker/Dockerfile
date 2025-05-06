# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o estimated .

# Final stage: minimal image
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/estimated .
COPY --from=builder /app/config.json .

# Expose port if needed
EXPOSE 8080

# Run the binary
ENTRYPOINT ["./estimated"]
