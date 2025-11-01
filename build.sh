#!/bin/bash

set -e

# Configuration
IMAGE_NAME="pocketstream-app-arm64-builder"
OUTPUT_BINARY="pocketstream-app-arm64"

echo "Building Docker image for ARM64..."
docker build --platform linux/arm64 -t $IMAGE_NAME .

echo "Extracting binary from Docker image..."
CONTAINER_ID=$(docker create --platform linux/arm64 $IMAGE_NAME)

if [ -d "./Pocketstream" ]; then
  rm -rf "./Pocketstream"
fi

mkdir ./Pocketstream

docker cp $CONTAINER_ID:/app/pocketstream-app ./Pocketstream/$OUTPUT_BINARY
docker rm $CONTAINER_ID
cp ./font.ttf ./Pocketstream/font.ttf

# Create muOS launcher script
cat > "./Pocketstream/mux_launch.sh" << 'EOF'
#!/bin/sh

cd /mnt/mmc/MUOS/application/Pocketstream
./pocketstream-app-arm64
EOF

chmod +x "./Pocketstream/mux_launch.sh"

echo "Build complete!"
echo "Binary: ./Pocketstream/$OUTPUT_BINARY"
echo "Binary info:"
file ./Pocketstream/$OUTPUT_BINARY
ls -lh ./Pocketstream/$OUTPUT_BINARY

