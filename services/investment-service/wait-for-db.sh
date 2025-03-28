#!/bin/sh
HOST=$(echo "$1" | cut -d ":" -f 1)
PORT=$(echo "$1" | cut -d ":" -f 2)
echo "Waiting for $HOST:$PORT..."
while ! nc -z $HOST $PORT; do
  echo "Database is unavailable - sleeping"
  sleep 1
done
echo "Database is up - starting service"
shift
exec "$@"
