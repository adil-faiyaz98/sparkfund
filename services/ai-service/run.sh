#!/bin/bash

# Build and run the Docker Compose setup
docker-compose up --build -d

# Wait for services to start
echo "Waiting for services to start..."
sleep 10

# Check if services are running
docker-compose ps

# Show logs
echo "Showing logs..."
docker-compose logs

echo "AI service is running at http://localhost:8000"
echo "API documentation is available at http://localhost:8000/docs"
echo "Health endpoint is available at http://localhost:8000/health"
