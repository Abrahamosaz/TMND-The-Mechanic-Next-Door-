#!/usr/bin/sh

BIN_DIR="./bin"
APP_BINARY="$BIN_DIR/app"

# Check if the binary exists
if [ ! -f "$APP_BINARY" ]; then
    echo "❌ Binary not found. Run './build.sh' first."
    exit 1
fi

echo "🚀 Running the application..."
"$APP_BINARY"