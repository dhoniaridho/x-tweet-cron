# Use official Go image
FROM golang:1.24-alpine

# Install necessary tools
RUN apk add --no-cache git curl

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the app
COPY . .

# Build the Go app
RUN go build -o app.exe

# Expose the callback port (optional, for OAuth redirect)
EXPOSE 8080

# Run the binary
CMD ["./app.exe"]
