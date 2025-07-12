#!/bin/sh

echo "🔍 DB_DRIVER=$DB_DRIVER"
echo "🔍 DB_SOURCE=$DB_SOURCE"

exec "$@"
# #!/bin/sh
# set -e

# echo "Waiting for Postgres..."

# # Простой цикл ожидания порта 5432
# until nc -z -v -w30 postgres 5432
# do
#   echo "Postgres is unavailable - sleeping"
#   sleep 1
# done

# echo "Postgres is up - running migrations..."

# /app/migrate -path ./migration -database "$DB_SOURCE" -verbose up

# echo "Starting app..."
# exec "$@"
