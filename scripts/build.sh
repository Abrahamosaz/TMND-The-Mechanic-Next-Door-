#!/usr/bin/sh

BIN_DIR="./bin"

# Create the bin directory if it doesn't exist
mkdir -p "$BIN_DIR"

# Build the Go project and place the binary in the bin directory
go build -o "$BIN_DIR/app" cmd/api/*.go

# Check if the build was successful
if [ $? -eq 0 ]; then
    echo "✅ Build successful. Run './scripts/app.sh' to start the app."
else
    echo "❌ Build failed. Please check for errors."
    exit 1
fi