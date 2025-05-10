FROM golang:1.24-alpine AS build

WORKDIR /app

# Install git and any other dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum* ./

# No need to modify go.mod version since we're using Go 1.24
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o rofl-service .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/rofl-service .

# Expose the application port
EXPOSE 8000

# Run the application
CMD ["./rofl-service"]