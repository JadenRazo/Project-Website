#!/bin/bash

# Database status check script
set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

# Load environment variables
if [ -f "$PROJECT_ROOT/backend/.env" ]; then
    export $(cat "$PROJECT_ROOT/backend/.env" | grep -v '^#' | xargs)
fi

# Database connection details
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-project_website}"
DB_USER="${DB_USER:-postgres}"

echo -e "${BLUE}Checking Database Status${NC}"
echo "Database: $DB_NAME on $DB_HOST:$DB_PORT"
echo ""

# Tables to check
REQUIRED_TABLES=(
    "shortened_urls"
    "url_clicks"
    "url_tags"
    "short_url_tags"
    "custom_domains"
    "incidents"
    "status_history"
    "word_filters"
)

# Check if we can connect
if ! PGPASSWORD="${DB_PASSWORD}" psql -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -U "$DB_USER" -c "SELECT 1;" >/dev/null 2>&1; then
    echo -e "${RED}❌ Cannot connect to database${NC}"
    echo "Please check your connection settings and try again."
    exit 1
fi

echo -e "${GREEN}✅ Database connection successful${NC}"
echo ""

# Check for each required table
echo -e "${BLUE}Checking for URL Shortener tables:${NC}"
missing_tables=()
found_tables=()

for table in "${REQUIRED_TABLES[@]}"; do
    if PGPASSWORD="${DB_PASSWORD}" psql -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -U "$DB_USER" -tc "SELECT 1 FROM information_schema.tables WHERE table_name = '$table';" 2>/dev/null | grep -q 1; then
        found_tables+=("$table")
        echo -e "  ${GREEN}✅${NC} $table"
    else
        missing_tables+=("$table")
        echo -e "  ${RED}❌${NC} $table (missing)"
    fi
done

echo ""

# Summary
if [ ${#missing_tables[@]} -eq 0 ]; then
    echo -e "${GREEN}🎉 Your database is fully updated!${NC}"
    echo "All URL shortener and monitoring tables are present."
else
    echo -e "${YELLOW}⚠️  Your database needs to be updated${NC}"
    echo ""
    echo "Missing tables: ${missing_tables[*]}"
    echo ""
    echo "To update your database, run:"
    echo -e "  ${BLUE}./scripts/update-db.sh --safe${NC}"
    echo ""
    echo "This will add the missing tables while preserving existing data."
fi