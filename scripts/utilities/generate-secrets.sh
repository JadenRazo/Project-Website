#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
ENV_FILE="$PROJECT_ROOT/.env"
ENV_EXAMPLE="$PROJECT_ROOT/.env.example"

echo "=================================="
echo "Portfolio Website Secrets Generator"
echo "=================================="
echo ""

if [ ! -f "$ENV_EXAMPLE" ]; then
    echo "ERROR: .env.example file not found at $ENV_EXAMPLE"
    exit 1
fi

if [ -f "$ENV_FILE" ]; then
    echo "Existing .env file found at: $ENV_FILE"
    echo ""
    read -p "Do you want to regenerate secrets? This will update existing secret values. (y/N): " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Operation cancelled."
        exit 0
    fi
    echo ""
    echo "Backing up existing .env to .env.backup..."
    cp "$ENV_FILE" "$ENV_FILE.backup"
    echo "Backup created: $ENV_FILE.backup"
else
    echo "No .env file found. Creating from .env.example..."
    cp "$ENV_EXAMPLE" "$ENV_FILE"
    echo "Created: $ENV_FILE"
fi

echo ""
echo "Generating secure secrets..."
echo ""

JWT_SECRET=$(openssl rand -base64 64 | tr -d '\n')
OAUTH_ENCRYPTION_KEY=$(openssl rand -hex 32)
DB_PASSWORD=$(openssl rand -base64 32 | tr -d '\n')
DB_ADMIN_PASSWORD=$(openssl rand -base64 32 | tr -d '\n')
DB_APP_PASSWORD=$(openssl rand -base64 32 | tr -d '\n')
DB_READONLY_PASSWORD=$(openssl rand -base64 32 | tr -d '\n')
ADMIN_TOKEN=$(openssl rand -base64 32 | tr -d '\n')
REDIS_APP_PASSWORD=$(openssl rand -base64 32 | tr -d '\n')
REDIS_MONITOR_PASSWORD=$(openssl rand -base64 32 | tr -d '\n')
REDIS_WORKER_PASSWORD=$(openssl rand -base64 32 | tr -d '\n')
REDIS_HEALTH_PASSWORD=$(openssl rand -base64 32 | tr -d '\n')
GRAFANA_PASSWORD=$(openssl rand -base64 32 | tr -d '\n')

update_or_append_env() {
    local key=$1
    local value=$2
    local file=$3

    if grep -q "^${key}=" "$file" 2>/dev/null; then
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' "s|^${key}=.*|${key}=${value}|" "$file"
        else
            sed -i "s|^${key}=.*|${key}=${value}|" "$file"
        fi
    else
        echo "${key}=${value}" >> "$file"
    fi
}

echo "Updating .env file with generated secrets..."

update_or_append_env "JWT_SECRET" "$JWT_SECRET" "$ENV_FILE"
update_or_append_env "OAUTH_ENCRYPTION_KEY" "$OAUTH_ENCRYPTION_KEY" "$ENV_FILE"
update_or_append_env "DB_PASSWORD" "$DB_PASSWORD" "$ENV_FILE"
update_or_append_env "DB_ADMIN_PASSWORD" "$DB_ADMIN_PASSWORD" "$ENV_FILE"
update_or_append_env "DB_APP_PASSWORD" "$DB_APP_PASSWORD" "$ENV_FILE"
update_or_append_env "DB_READONLY_PASSWORD" "$DB_READONLY_PASSWORD" "$ENV_FILE"
update_or_append_env "ADMIN_TOKEN" "$ADMIN_TOKEN" "$ENV_FILE"
update_or_append_env "REDIS_APP_PASSWORD" "$REDIS_APP_PASSWORD" "$ENV_FILE"
update_or_append_env "REDIS_MONITOR_PASSWORD" "$REDIS_MONITOR_PASSWORD" "$ENV_FILE"
update_or_append_env "REDIS_WORKER_PASSWORD" "$REDIS_WORKER_PASSWORD" "$ENV_FILE"
update_or_append_env "REDIS_HEALTH_PASSWORD" "$REDIS_HEALTH_PASSWORD" "$ENV_FILE"
update_or_append_env "GRAFANA_PASSWORD" "$GRAFANA_PASSWORD" "$ENV_FILE"

echo ""
echo "✓ Secrets generated and updated successfully!"
echo ""
echo "Generated secrets:"
echo "  - JWT_SECRET"
echo "  - OAUTH_ENCRYPTION_KEY"
echo "  - DB_PASSWORD (legacy, kept for compatibility)"
echo "  - DB_ADMIN_PASSWORD (PostgreSQL superuser)"
echo "  - DB_APP_PASSWORD (PostgreSQL application user)"
echo "  - DB_READONLY_PASSWORD (PostgreSQL read-only user)"
echo "  - ADMIN_TOKEN"
echo "  - REDIS_APP_PASSWORD"
echo "  - REDIS_MONITOR_PASSWORD"
echo "  - REDIS_WORKER_PASSWORD"
echo "  - REDIS_HEALTH_PASSWORD"
echo "  - GRAFANA_PASSWORD"
echo ""
echo "IMPORTANT: You still need to configure:"
echo "  - OAuth provider credentials (GOOGLE_CLIENT_ID, GITHUB_CLIENT_ID, MICROSOFT_CLIENT_ID, etc.)"
echo "  - Admin allowed emails (ADMIN_ALLOWED_EMAILS)"
echo "  - Database connection details if different from defaults (DB_HOST, DB_NAME)"
echo ""

# Generate PostgreSQL SSL/TLS certificates
echo "Generating PostgreSQL SSL/TLS certificates..."
if [ -f "$SCRIPT_DIR/generate-postgres-certs.sh" ]; then
    "$SCRIPT_DIR/generate-postgres-certs.sh"
else
    echo "  Warning: generate-postgres-certs.sh not found, skipping certificate generation"
    echo "  Run manually: ./scripts/utilities/generate-postgres-certs.sh"
fi

echo ""
echo "Your .env file is ready at: $ENV_FILE"

if [ -f "$ENV_FILE.backup" ]; then
    echo "Previous configuration backed up to: $ENV_FILE.backup"
fi

echo ""
echo "=================================="
echo "Next steps:"
echo "1. Review and update .env file with your specific configuration"
echo "2. Configure OAuth credentials from provider consoles"
echo "3. Update ADMIN_ALLOWED_EMAILS with your email address"
echo "4. Run: docker-compose up -d"
echo "=================================="
