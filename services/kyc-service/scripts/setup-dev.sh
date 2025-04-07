#!/bin/bash

# Create necessary directories
mkdir -p tmp
mkdir -p data/documents

# Copy environment file if it doesn't exist
if [ ! -f .env ]; then
    cp .env.example .env
    echo "Created .env file from .env.example"
fi

# Download ML models if they don't exist
if [ ! -d "internal/ai-powered/models" ]; then
    mkdir -p internal/ai-powered/models
    echo "Created models directory. Please download required ML models."
fi

# Start development environment
docker-compose -f docker-compose.dev.yml up --build