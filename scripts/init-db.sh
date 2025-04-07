#!/bin/bash
set -e

# Create databases
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE kyc;
    CREATE DATABASE investment;
    CREATE DATABASE user_service;
EOSQL

# Create users and grant privileges
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER kyc WITH PASSWORD 'kyc';
    GRANT ALL PRIVILEGES ON DATABASE kyc TO kyc;
    
    CREATE USER investment WITH PASSWORD 'investment';
    GRANT ALL PRIVILEGES ON DATABASE investment TO investment;
    
    CREATE USER user_service WITH PASSWORD 'user';
    GRANT ALL PRIVILEGES ON DATABASE user_service TO user_service;
EOSQL

# Connect to each database and create extensions
for DB in kyc investment user_service; do
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$DB" <<-EOSQL
        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
        CREATE EXTENSION IF NOT EXISTS "pgcrypto";
        CREATE EXTENSION IF NOT EXISTS "pg_trgm";
EOSQL
done

echo "Database initialization completed"
