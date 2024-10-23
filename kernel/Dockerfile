FROM golang:1.23

WORKDIR /kernel

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Copy the source code
COPY js.go .
COPY main.go .
COPY ./kernel ./kernel

# Build the WASM module
RUN GOOS=js GOARCH=wasm go build -o space-wars.wasm

