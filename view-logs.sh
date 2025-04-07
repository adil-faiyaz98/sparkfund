#!/bin/bash
echo "==================================================="
echo "SparkFund - Viewing Service Logs"
echo "==================================================="

if [ -z "$1" ]; then
    echo "Viewing logs for all services..."
    docker-compose logs -f
else
    echo "Viewing logs for $1..."
    docker-compose logs -f $1
fi

echo "==================================================="
