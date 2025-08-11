#!/bin/bash

# Database migration script

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | xargs)
fi

# Default values
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-password}
DB_NAME=${DB_NAME:-crud_demo}

# Database URL
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

case "$1" in
    "up")
        echo "Running database migrations..."
        # In a real project, you would use a migration tool like golang-migrate
        # For now, we'll use psql to run SQL files
        if command -v psql >/dev/null 2>&1; then
            psql "$DATABASE_URL" -f internal/database/migrations/001_create_users_table.sql
            echo "Migrations completed successfully"
        else
            echo "psql not found. Using Go application to run migrations..."
            go run cmd/migrate/main.go up
        fi
        ;;
    "down")
        echo "Rolling back database migrations..."
        if command -v psql >/dev/null 2>&1; then
            psql "$DATABASE_URL" -c "DROP TABLE IF EXISTS users CASCADE;"
            echo "Migrations rolled back successfully"
        else
            echo "psql not found. Using Go application to rollback migrations..."
            go run cmd/migrate/main.go down
        fi
        ;;
    "reset")
        echo "Resetting database..."
        $0 down
        $0 up
        ;;
    "status")
        echo "Checking migration status..."
        if command -v psql >/dev/null 2>&1; then
            psql "$DATABASE_URL" -c "\dt"
        else
            echo "psql not found. Cannot check migration status"
        fi
        ;;
    *)
        echo "Usage: $0 {up|down|reset|status}"
        echo "  up     - Run pending migrations"
        echo "  down   - Rollback migrations"
        echo "  reset  - Rollback and re-run all migrations"
        echo "  status - Show current migration status"
        exit 1
        ;;
esac