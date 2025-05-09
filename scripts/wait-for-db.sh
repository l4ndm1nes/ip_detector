#!/bin/sh
set -e

host="$1"
shift
cmd="$@"

until pg_isready -h "$host" -p 5432 > /dev/null 2>&1; do
  echo "Waiting for postgres at $host..."
  sleep 2
done

echo "Postgres is ready, starting the application..."
exec $cmd
