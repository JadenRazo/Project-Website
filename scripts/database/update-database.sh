#!/bin/bash

# =============================================================================
# Database Update Script for Portfolio Website
# =============================================================================
# This script safely updates your PostgreSQL database with the new schema
# =============================================================================

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"
BACKEND_DIR="$PROJECT_ROOT/backend"
SCHEMA_FILE="$BACKEND_DIR/schema.sql"
BACKUP_DIR="$PROJECT_ROOT/database-backups"

# Default database configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-project_website}"
DB_USER="${DB_USER:-postgres}"

log() {
    echo -e "${BLUE}[$(date +'%H:%M:%S')] $1${NC}"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

# Load environment variables from backend/.env if it exists
load_env() {
    if [ -f "$BACKEND_DIR/.env" ]; then
        log "Loading environment variables from backend/.env"
        set -a
        source "$BACKEND_DIR/.env"
        set +a
        
        # Update database configuration from env
        DB_HOST="${DB_HOST:-localhost}"
        DB_PORT="${DB_PORT:-5432}"
        DB_NAME="${DB_NAME:-project_website}"
        DB_USER="${DB_USER:-postgres}"
    fi
}

# Check if PostgreSQL client is installed
check_dependencies() {
    log "Checking dependencies..."
    
    if ! command -v psql &> /dev/null; then
        error "PostgreSQL client (psql) is required but not installed. Please install it first."
    fi
    
    if ! command -v pg_dump &> /dev/null; then
        error "pg_dump is required but not installed. Please install PostgreSQL client tools."
    fi
    
    success "All dependencies found"
}

# Test database connection
test_connection() {
    log "Testing database connection..."
    
    if ! PGPASSWORD="${DB_PASSWORD}" psql -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -U "$DB_USER" -c "SELECT 1;" >/dev/null 2>&1; then
        error "Cannot connect to database. Please check your connection settings:
  Host: $DB_HOST
  Port: $DB_PORT
  Database: $DB_NAME
  User: $DB_USER"
    fi
    
    success "Database connection successful"
}

# Create backup directory if it doesn't exist
create_backup_dir() {
    if [ ! -d "$BACKUP_DIR" ]; then
        log "Creating backup directory..."
        mkdir -p "$BACKUP_DIR"
    fi
}

# Backup the current database
backup_database() {
    log "Creating database backup..."
    
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local backup_file="$BACKUP_DIR/backup_${DB_NAME}_${timestamp}.sql"
    
    if PGPASSWORD="${DB_PASSWORD}" pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME" > "$backup_file"; then
        success "Database backed up to: $backup_file"
        echo "$backup_file"  # Return the backup file path
    else
        error "Failed to create database backup"
    fi
}

# Apply the schema updates
apply_schema() {
    log "Applying schema updates..."
    
    if [ ! -f "$SCHEMA_FILE" ]; then
        error "Schema file not found: $SCHEMA_FILE"
    fi
    
    # Create a temporary file with transaction wrapper
    local temp_file=$(mktemp)
    cat > "$temp_file" << EOF
-- Start transaction for safe schema update
BEGIN;

-- Set client encoding
SET client_encoding = 'UTF8';

-- Include the schema file
\i $SCHEMA_FILE

-- If we get here, everything worked
COMMIT;
EOF
    
    # Apply the schema
    if PGPASSWORD="${DB_PASSWORD}" psql -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -U "$DB_USER" -f "$temp_file" -v ON_ERROR_STOP=1; then
        success "Schema updates applied successfully"
        rm -f "$temp_file"
    else
        error "Failed to apply schema updates. The database has been rolled back to its previous state."
        rm -f "$temp_file"
    fi
}

# Verify the update
verify_update() {
    log "Verifying database update..."
    
    # Check for key tables
    local tables=("users" "shortened_urls" "url_clicks" "incidents" "status_history" "messaging_channels" "word_filters")
    local missing_tables=()
    
    for table in "${tables[@]}"; do
        if ! PGPASSWORD="${DB_PASSWORD}" psql -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -U "$DB_USER" -tc "SELECT 1 FROM information_schema.tables WHERE table_name = '$table';" | grep -q 1; then
            missing_tables+=("$table")
        fi
    done
    
    if [ ${#missing_tables[@]} -eq 0 ]; then
        success "All expected tables are present"
    else
        warning "The following tables are missing: ${missing_tables[*]}"
    fi
    
    # Show table count
    local table_count=$(PGPASSWORD="${DB_PASSWORD}" psql -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -U "$DB_USER" -tc "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';")
    log "Total tables in database: $table_count"
}

# Show update summary
show_summary() {
    echo ""
    echo -e "${GREEN}🎉 Database Update Complete!${NC}"
    echo ""
    echo -e "${CYAN}Summary:${NC}"
    echo "  • Database: $DB_NAME"
    echo "  • Schema file: $SCHEMA_FILE"
    echo "  • Backup created: $1"
    echo ""
    echo -e "${CYAN}Key tables added/updated:${NC}"
    echo "  • URL Shortener: shortened_urls, url_clicks, url_tags"
    echo "  • Status Monitoring: incidents, status_history"
    echo "  • Messaging: word_filters"
    echo "  • All supporting tables and indexes"
    echo ""
    echo -e "${YELLOW}Next steps:${NC}"
    echo "  1. Start the development environment: ./start-dev.sh"
    echo "  2. The URL shortener will be available at http://localhost:3000/urlshortener"
    echo "  3. API endpoints are at http://localhost:8083"
    echo ""
}

# Main execution
main() {
    echo -e "${CYAN}╔══════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║             PostgreSQL Database Update Script                     ║${NC}"
    echo -e "${CYAN}╚══════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    
    # Load environment variables
    load_env
    
    # Run checks
    check_dependencies
    test_connection
    
    # Confirm before proceeding
    echo ""
    warning "This will update your database schema. A backup will be created first."
    echo -e "${YELLOW}Database: $DB_NAME on $DB_HOST:$DB_PORT${NC}"
    echo ""
    read -p "Do you want to proceed? (y/N) " -n 1 -r
    echo ""
    
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log "Update cancelled"
        exit 0
    fi
    
    # Create backup
    create_backup_dir
    backup_file=$(backup_database)
    
    # Apply schema
    apply_schema
    
    # Verify
    verify_update
    
    # Show summary
    show_summary "$backup_file"
}

# Handle script interruption
trap 'echo -e "\n${YELLOW}Update interrupted.${NC}"; exit 1' INT

# Run main function
main