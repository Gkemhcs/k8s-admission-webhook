# Use a minimal base image with Go support
FROM golang:1.23-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the webhook binary
RUN go build -o webhook-server main.go

# Use a lightweight final image
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Install required certificates
RUN apk --no-cache add ca-certificates

# Copy the compiled binary from builder stage
COPY --from=builder /app/webhook-server .

# Expose the webhook server port
EXPOSE 8080

# Command to run the webhook server with certs passed as flags
ENTRYPOINT ["./webhook-server"] 
