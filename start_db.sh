#!/usr/bin/env bash

set -e

cd ~/development/callisto

docker stop callisto-db || true
docker kill callisto-db || true
docker rm callisto-db || true

docker run \
  --name callisto-db \
  -e POSTGRES_USER=callisto \
  -e POSTGRES_PASSWORD=callisto \
  -e POSTGRES_DB=callisto \
  -v "$(pwd)"/database:/database \
  -p 5432:5432 \
  -d \
  postgres:latest

sleep 2

SQL_FILES=$(find database/schema -type f -name "*.sql" | xargs -I{} basename {} | sort)

for file in $SQL_FILES; do
  echo "Running $file"
  docker exec -i callisto-db psql -U callisto -d callisto -f /database/schema/$file 1>/dev/null
done

echo "Database started and schema loaded"
