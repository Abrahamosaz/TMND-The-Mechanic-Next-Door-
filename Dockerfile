# Use the Go 1.23 alpine official image
# https://hub.docker.com/_/golang
FROM golang:1.23-alpine

# Create and change to the app directory.
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Copy local code to the container image.
COPY . .

# Install project dependencies
RUN go mod download

# Build the application using the script
RUN go build -o "./bin/app" cmd/api/*.go

# Run the service on container startup
ENTRYPOINT ["./bin/app"]