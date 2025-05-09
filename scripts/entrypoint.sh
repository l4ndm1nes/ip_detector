#!/bin/sh
set -e

echo "Waiting for PostgreSQL to be ready..."
until pg_isready -h db -p 5432 > /dev/null 2>&1; do
  sleep 2
done

echo "PostgreSQL is ready. Applying migrations..."
/app/main migrate -path=./migrations -database="postgres://${DB_USER}:${DB_PASSWORD}@db:5432/${DB_NAME}?sslmode=disable" up

echo "Starting application..."
exec "$@"
