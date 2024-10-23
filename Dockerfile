# Use the official Golang image as the base image
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the WASM module
RUN GOOS=js GOARCH=wasm go build -o ./client/public/space-wars.wasm

# Set the entrypoint to a no-op command
CMD ["/bin/sh", "-c", "echo 'WASM module built'"]