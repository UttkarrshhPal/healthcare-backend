#!/bin/bash

echo "Setting up Healthcare Portal..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Creating .env file from .env.example..."
    cp .env.example .env
    echo "Please update .env file with your database credentials"
    exit 1
fi

# Load environment variables
source .env

# Check if PostgreSQL is running
echo "Checking PostgreSQL connection..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d postgres -c '\q' 2>/dev/null
if [ $? -ne 0 ]; then
    echo "Error: Cannot connect to PostgreSQL. Please ensure PostgreSQL is running and credentials are correct."
    exit 1
fi

# Create database if it doesn't exist
echo "Creating database if not exists..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME;" 2>/dev/null

# Run the application to trigger migrations
echo "Running database migrations..."
go run cmd/server/main.go &
SERVER_PID=$!
sleep 5
kill $SERVER_PID 2>/dev/null

# Run seeder
echo "Seeding database with default users..."
go run cmd/seeder/main.go

echo "Setup completed!"