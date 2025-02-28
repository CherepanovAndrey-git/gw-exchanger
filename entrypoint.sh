#!/bin/sh
set -e

echo "Waiting for the exchanger database to be ready..."
/app/wait-for-db.sh exchanger-db

echo "Running migrations..."
/app/goose -dir /app/sql/schema postgres "$DB_URL" up

echo "Starting gRPC server..."
/app/main