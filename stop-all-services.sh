#!/bin/bash
echo "==================================================="
echo "SparkFund - Stopping All Microservices"
echo "==================================================="

echo "Stopping and removing all containers..."
docker-compose -f docker-compose-all.yml down

echo "==================================================="
echo "All services have been stopped."
echo "==================================================="
