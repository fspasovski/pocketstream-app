FROM arm64v8/ubuntu:22.04

# Prevent interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install Go, SDL2, and build dependencies
RUN apt-get update && apt-get install -y \
    wget \
    gcc \
    g++ \
    make \
    pkg-config \
    libsdl2-dev \
    libsdl2-image-dev \
    libsdl2-mixer-dev \
    libsdl2-ttf-dev \
    libsdl2-gfx-dev \
    && rm -rf /var/lib/apt/lists/*

# Install Go (use a recent stable version)
RUN wget -q https://go.dev/dl/go1.23.2.linux-arm64.tar.gz && \
    tar -C /usr/local -xzf go1.23.2.linux-arm64.tar.gz && \
    rm go1.23.2.linux-arm64.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=arm64

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum* ./
RUN go mod download || true

# Copy the rest of the application
COPY . .

# Build the application with optimizations
RUN go build -ldflags="-s -w" -o pocketstream-app .

# The binary will be at /app/myapp