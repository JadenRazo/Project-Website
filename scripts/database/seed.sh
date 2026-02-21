#!/bin/bash

[ -f .env ] && export $(grep -v '^#' .env | xargs)

# Color codes for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Database configuration
DB_NAME="${DB_NAME:-portfolio_db}"
DB_USER="${DB_USER:-postgres}"
DB_HOST="${DB_HOST:-localhost}"

echo -e "${GREEN}Starting database seeding...${NC}"

# Run the seed SQL file
psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f "backend/db/seed.sql"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Database seeded successfully!${NC}"
else
    echo -e "${RED}Database seeding failed!${NC}"
    exit 1
fi
