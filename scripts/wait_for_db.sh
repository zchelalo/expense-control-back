#!/bin/bash

DB_HOST=$1
DB_PORT=$2

if [ -z "$DB_PORT" ]; then
  DB_PORT="5432"
fi

until nc -z -v -w30 $DB_HOST $DB_PORT
do
  echo "Waiting for database at $DB_HOST:$DB_PORT..."
  sleep 1
done

echo "Database is up and running at $DB_HOST:$DB_PORT"