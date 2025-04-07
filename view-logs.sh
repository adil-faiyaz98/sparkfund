#!/bin/bash
echo "==================================================="
echo "SparkFund - Viewing Service Logs"
echo "==================================================="

if [ -z "$1" ]; then
    echo "Viewing logs for all services..."
    docker-compose -f docker-compose-all.yml logs -f
else
    echo "Viewing logs for $1..."
    docker-compose -f docker-compose-all.yml logs -f $1
fi

echo "==================================================="
