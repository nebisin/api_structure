#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path ./migrations -database "$DB_URI" -verbose up

echo "start the app"
exec "$@"