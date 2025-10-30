#!/bin/bash

# Simple database update wrapper
set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}PostgreSQL Database Update${NC}"
echo ""

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

# Use the main schema file
echo -e "${YELLOW}Using database schema from backend/schema.sql...${NC}"
SCHEMA_FILE="$PROJECT_ROOT/backend/schema.sql"

if [ ! -f "$SCHEMA_FILE" ]; then
    echo -e "${RED}Error: Schema file not found at $SCHEMA_FILE${NC}"
    exit 1
fi

# Load environment variables
if [ -f "$PROJECT_ROOT/backend/.env" ]; then
    export $(cat "$PROJECT_ROOT/backend/.env" | grep -v '^#' | xargs)
fi

# Database connection details
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-project_website}"
DB_USER="${DB_USER:-postgres}"

echo "Database: $DB_NAME on $DB_HOST:$DB_PORT"
echo ""

# Apply the schema
echo -e "${BLUE}Applying database schema...${NC}"
PGPASSWORD="${DB_PASSWORD}" psql -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -U "$DB_USER" -f "$SCHEMA_FILE"

echo ""
echo -e "${GREEN}✅ Database schema applied successfully!${NC}"
echo ""
echo "All tables have been created/updated in your database."
echo "You can now start the development environment with: ./start-dev.sh"