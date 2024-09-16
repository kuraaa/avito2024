#!/bin/sh

until pg_isready -h avitoDb -p 5432 -U user; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 5
done

echo "PostgreSQL is ready!"
