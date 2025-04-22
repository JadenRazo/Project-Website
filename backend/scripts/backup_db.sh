#!/bin/bash
set -e

# Print colorful messages
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}Project-Website Database Backup Script${NC}"
echo -e "${BLUE}========================================${NC}"

# Navigate to backend directory
cd "$(dirname "${BASH_SOURCE[0]}")/.."
BACKEND_DIR="$(pwd)"
echo -e "${GREEN}Backend directory:${NC} ${YELLOW}${BACKEND_DIR}${NC}"

# Set default variables
DB_TYPE=${DB_TYPE:-"sqlite"}
BACKUP_DIR="${BACKEND_DIR}/backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/backup_${TIMESTAMP}.sql"

# Create backup directory if it doesn't exist
mkdir -p "${BACKUP_DIR}"

# Load config values if possible
if [ -f "${BACKEND_DIR}/config/development.yaml" ]; then
    echo -e "${GREEN}Loading database configuration from development.yaml...${NC}"
    # This is a simple extraction - for production use a proper YAML parser
    if grep -q "database:" "${BACKEND_DIR}/config/development.yaml"; then
        DB_TYPE=$(grep -A 10 "database:" "${BACKEND_DIR}/config/development.yaml" | grep "driver:" | head -n 1 | awk -F': ' '{print $2}' | tr -d ' "')
        DB_DSN=$(grep -A 10 "database:" "${BACKEND_DIR}/config/development.yaml" | grep "dsn:" | head -n 1 | awk -F': ' '{print $2}' | tr -d ' "')
        
        # Extract database credentials if using PostgreSQL or MySQL
        if [[ "$DB_TYPE" == "postgres" ]]; then
            DB_HOST=$(echo $DB_DSN | grep -o "host=[^ ]*" | cut -d= -f2)
            DB_PORT=$(echo $DB_DSN | grep -o "port=[^ ]*" | cut -d= -f2)
            DB_NAME=$(echo $DB_DSN | grep -o "dbname=[^ ]*" | cut -d= -f2)
            DB_USER=$(echo $DB_DSN | grep -o "user=[^ ]*" | cut -d= -f2)
            DB_PASSWORD=$(echo $DB_DSN | grep -o "password=[^ ]*" | cut -d= -f2)
        elif [[ "$DB_TYPE" == "mysql" ]]; then
            DB_USER=$(echo $DB_DSN | cut -d: -f1)
            DB_PASSWORD=$(echo $DB_DSN | cut -d: -f2 | cut -d@ -f1)
            DB_HOST=$(echo $DB_DSN | cut -d@ -f2 | cut -d/ -f1 | cut -d: -f1)
            DB_PORT=$(echo $DB_DSN | cut -d@ -f2 | cut -d/ -f1 | cut -d: -f2)
            DB_NAME=$(echo $DB_DSN | cut -d/ -f2 | cut -d? -f1)
        elif [[ "$DB_TYPE" == "sqlite" ]]; then
            DB_FILE=$DB_DSN
        fi
    fi
fi

echo -e "\n${GREEN}Database backup information:${NC}"
echo -e "  Database type: ${YELLOW}${DB_TYPE}${NC}"
echo -e "  Backup file: ${YELLOW}${BACKUP_FILE}${NC}"

# Perform backup based on database type
case $DB_TYPE in
    "sqlite")
        if [ -z "$DB_FILE" ]; then
            DB_FILE="${BACKEND_DIR}/data/database.db"
            echo -e "${YELLOW}Using default SQLite database location: ${DB_FILE}${NC}"
        fi
        
        if [ ! -f "$DB_FILE" ]; then
            echo -e "${RED}Error: SQLite database file not found: ${DB_FILE}${NC}"
            exit 1
        fi
        
        echo -e "\n${GREEN}Backing up SQLite database...${NC}"
        if command -v sqlite3 &> /dev/null; then
            sqlite3 "$DB_FILE" ".dump" > "${BACKUP_FILE}"
            # Also create a binary copy as .db file
            cp "$DB_FILE" "${BACKUP_DIR}/backup_${TIMESTAMP}.db"
            echo -e "${GREEN}SQLite database backup completed successfully!${NC}"
            echo -e "  SQL Dump: ${YELLOW}${BACKUP_FILE}${NC}"
            echo -e "  Binary copy: ${YELLOW}${BACKUP_DIR}/backup_${TIMESTAMP}.db${NC}"
        else
            echo -e "${RED}Error: sqlite3 command not found. Please install SQLite.${NC}"
            exit 1
        fi
        ;;
        
    "postgres")
        if [ -z "$DB_HOST" ] || [ -z "$DB_NAME" ]; then
            echo -e "${RED}Error: Missing PostgreSQL connection details.${NC}"
            exit 1
        fi
        
        echo -e "\n${GREEN}Backing up PostgreSQL database...${NC}"
        if command -v pg_dump &> /dev/null; then
            PGPASSWORD=$DB_PASSWORD pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -F c -f "${BACKUP_FILE}.pgdump"
            echo -e "${GREEN}PostgreSQL database backup completed successfully!${NC}"
            echo -e "  Backup file: ${YELLOW}${BACKUP_FILE}.pgdump${NC}"
        else
            echo -e "${RED}Error: pg_dump command not found. Please install PostgreSQL client tools.${NC}"
            exit 1
        fi
        ;;
        
    "mysql")
        if [ -z "$DB_HOST" ] || [ -z "$DB_NAME" ]; then
            echo -e "${RED}Error: Missing MySQL connection details.${NC}"
            exit 1
        fi
        
        echo -e "\n${GREEN}Backing up MySQL database...${NC}"
        if command -v mysqldump &> /dev/null; then
            mysqldump -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD $DB_NAME > "${BACKUP_FILE}"
            echo -e "${GREEN}MySQL database backup completed successfully!${NC}"
            echo -e "  Backup file: ${YELLOW}${BACKUP_FILE}${NC}"
        else
            echo -e "${RED}Error: mysqldump command not found. Please install MySQL client tools.${NC}"
            exit 1
        fi
        ;;
        
    *)
        echo -e "${RED}Error: Unsupported database type: ${DB_TYPE}${NC}"
        echo -e "  Supported types: sqlite, postgres, mysql"
        exit 1
        ;;
esac

# Create a list of backups
echo -e "\n${GREEN}Available backups:${NC}"
ls -lh "${BACKUP_DIR}" | grep -v "total" | sort -r

echo -e "\n${GREEN}Backup operation completed successfully!${NC}"
