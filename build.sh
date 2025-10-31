#!/bin/bash

set -e

# Configuration
IMAGE_NAME="pocketstream-app-arm64-builder"
OUTPUT_BINARY="pocketstream-app-arm64"

echo "Building Docker image for ARM64..."
docker build --platform linux/arm64 -t $IMAGE_NAME .

echo "Extracting binary from Docker image..."
CONTAINER_ID=$(docker create --platform linux/arm64 $IMAGE_NAME)
docker cp $CONTAINER_ID:/app/pocketstream-app ./$OUTPUT_BINARY
docker rm $CONTAINER_ID

echo "Build complete!"
echo "Binary: ./$OUTPUT_BINARY"
echo "Binary info:"
file $OUTPUT_BINARY
ls -lh $OUTPUT_BINARY

echo "To verify the binary, you can run:"
echo "   ldd $OUTPUT_BINARY  # (on Linux system) to check dependencies"
echo "Transfer to your device and run!"