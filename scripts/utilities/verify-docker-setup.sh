#!/bin/bash

# Docker Setup Verification Script
# Checks if all required files and configurations are in place

set -e

PROJECT_ROOT="/main/Project-Website"
ERRORS=0
WARNINGS=0

# Color codes
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo "======================================"
echo "Docker Setup Verification"
echo "======================================"
echo ""

# Function to check file exists
check_file() {
    local file=$1
    local description=$2
    if [ -f "$file" ]; then
        echo -e "${GREEN}✓${NC} $description: $file"
    else
        echo -e "${RED}✗${NC} $description: $file (MISSING)"
        ((ERRORS++))
    fi
}

# Function to check directory exists
check_dir() {
    local dir=$1
    local description=$2
    if [ -d "$dir" ]; then
        echo -e "${GREEN}✓${NC} $description: $dir"
    else
        echo -e "${RED}✗${NC} $description: $dir (MISSING)"
        ((ERRORS++))
    fi
}

# Function to check file contains string
check_content() {
    local file=$1
    local search=$2
    local description=$3
    if grep -q "$search" "$file" 2>/dev/null; then
        echo -e "${GREEN}✓${NC} $description"
    else
        echo -e "${YELLOW}⚠${NC} $description (NOT FOUND)"
        ((WARNINGS++))
    fi
}

echo "Checking Dockerfiles..."
echo "------------------------"
check_file "$PROJECT_ROOT/frontend/Dockerfile" "Frontend Dockerfile"
check_file "$PROJECT_ROOT/frontend/nginx.conf" "Frontend nginx config"
check_file "$PROJECT_ROOT/backend/deployments/docker/api/Dockerfile" "API Dockerfile"
check_file "$PROJECT_ROOT/backend/deployments/docker/devpanel/Dockerfile" "DevPanel Dockerfile"
check_file "$PROJECT_ROOT/backend/deployments/docker/messaging/Dockerfile" "Messaging Dockerfile"
check_file "$PROJECT_ROOT/backend/deployments/docker/urlshortener/Dockerfile" "URL Shortener Dockerfile"
check_file "$PROJECT_ROOT/backend/deployments/docker/worker/Dockerfile" "Worker Dockerfile"
echo ""

echo "Checking .dockerignore files..."
echo "--------------------------------"
check_file "$PROJECT_ROOT/frontend/.dockerignore" "Frontend .dockerignore"
check_file "$PROJECT_ROOT/backend/.dockerignore" "Backend .dockerignore"
echo ""

echo "Checking Configuration Files..."
echo "--------------------------------"
check_file "$PROJECT_ROOT/docker-compose.yml" "docker-compose.yml"
check_file "$PROJECT_ROOT/.env.example" ".env.example"
check_file "$PROJECT_ROOT/deploy/redis/redis.conf" "Redis configuration"
check_file "$PROJECT_ROOT/deploy/redis/users.acl" "Redis ACL users"
check_file "$PROJECT_ROOT/deploy/prometheus/prometheus.yml" "Prometheus configuration"
echo ""

echo "Checking Redis TLS Certificates..."
echo "-----------------------------------"
check_dir "$PROJECT_ROOT/deploy/redis/certs" "Certificate directory"
check_file "$PROJECT_ROOT/deploy/redis/certs/ca.crt" "CA certificate"
check_file "$PROJECT_ROOT/deploy/redis/certs/ca.key" "CA private key"
check_file "$PROJECT_ROOT/deploy/redis/certs/redis.crt" "Redis certificate"
check_file "$PROJECT_ROOT/deploy/redis/certs/redis.key" "Redis private key"
echo ""

echo "Checking Documentation..."
echo "-------------------------"
check_file "$PROJECT_ROOT/DOCKER.md" "Docker deployment guide"
check_file "$PROJECT_ROOT/DEPLOYMENT_CHECKLIST.md" "Deployment checklist"
check_file "$PROJECT_ROOT/scripts/utilities/generate-redis-certs.sh" "Certificate generation script"
echo ""

echo "Checking docker-compose.yml structure..."
echo "-----------------------------------------"
check_content "$PROJECT_ROOT/docker-compose.yml" "  api:" "API service defined"
check_content "$PROJECT_ROOT/docker-compose.yml" "  devpanel:" "DevPanel service defined"
check_content "$PROJECT_ROOT/docker-compose.yml" "  messaging:" "Messaging service defined"
check_content "$PROJECT_ROOT/docker-compose.yml" "  urlshortener:" "URL Shortener service defined"
check_content "$PROJECT_ROOT/docker-compose.yml" "  worker:" "Worker service defined"
check_content "$PROJECT_ROOT/docker-compose.yml" "8080:8080" "API port mapping"
check_content "$PROJECT_ROOT/docker-compose.yml" "8081:8081" "DevPanel port mapping"
check_content "$PROJECT_ROOT/docker-compose.yml" "8082:8082" "Messaging port mapping"
check_content "$PROJECT_ROOT/docker-compose.yml" "8083:8083" "URL Shortener port mapping"
check_content "$PROJECT_ROOT/docker-compose.yml" "8084:8084" "Worker port mapping"
echo ""

echo "Checking .env.example variables..."
echo "----------------------------------"
check_content "$PROJECT_ROOT/.env.example" "REDIS_APP_PASSWORD" "REDIS_APP_PASSWORD"
check_content "$PROJECT_ROOT/.env.example" "REDIS_MONITOR_PASSWORD" "REDIS_MONITOR_PASSWORD"
check_content "$PROJECT_ROOT/.env.example" "REDIS_WORKER_PASSWORD" "REDIS_WORKER_PASSWORD"
check_content "$PROJECT_ROOT/.env.example" "REDIS_HEALTH_PASSWORD" "REDIS_HEALTH_PASSWORD"
check_content "$PROJECT_ROOT/.env.example" "ADMIN_TOKEN" "ADMIN_TOKEN"
check_content "$PROJECT_ROOT/.env.example" "GRAFANA_PASSWORD" "GRAFANA_PASSWORD"
check_content "$PROJECT_ROOT/.env.example" "REDIS_TLS_ENABLED" "REDIS_TLS_ENABLED"
echo ""

echo "Checking Redis TLS configuration..."
echo "------------------------------------"
check_content "$PROJECT_ROOT/deploy/redis/redis.conf" "tls-port 6379" "TLS port enabled"
check_content "$PROJECT_ROOT/deploy/redis/redis.conf" "TLSv1.3" "TLS 1.3 configured"
check_content "$PROJECT_ROOT/deploy/redis/redis.conf" "tls-cert-file" "TLS certificate path"
check_content "$PROJECT_ROOT/deploy/redis/redis.conf" "tls-key-file" "TLS key path"
echo ""

echo "Checking Prometheus scrape targets..."
echo "--------------------------------------"
check_content "$PROJECT_ROOT/deploy/prometheus/prometheus.yml" "api:8080" "API scrape target"
check_content "$PROJECT_ROOT/deploy/prometheus/prometheus.yml" "devpanel:8081" "DevPanel scrape target"
check_content "$PROJECT_ROOT/deploy/prometheus/prometheus.yml" "messaging:8082" "Messaging scrape target"
check_content "$PROJECT_ROOT/deploy/prometheus/prometheus.yml" "urlshortener:8083" "URL Shortener scrape target"
check_content "$PROJECT_ROOT/deploy/prometheus/prometheus.yml" "worker:8084" "Worker scrape target"
echo ""

echo "Checking environment file..."
echo "----------------------------"
if [ -f "$PROJECT_ROOT/.env" ]; then
    echo -e "${GREEN}✓${NC} .env file exists"

    # Check for placeholder values
    if grep -q "your_.*_password_here" "$PROJECT_ROOT/.env" 2>/dev/null; then
        echo -e "${YELLOW}⚠${NC} .env contains placeholder passwords - update before deployment"
        ((WARNINGS++))
    fi

    if grep -q "your_.*_token_here" "$PROJECT_ROOT/.env" 2>/dev/null; then
        echo -e "${YELLOW}⚠${NC} .env contains placeholder tokens - update before deployment"
        ((WARNINGS++))
    fi
else
    echo -e "${YELLOW}⚠${NC} .env file not found - copy from .env.example"
    echo "   Run: cp .env.example .env"
    ((WARNINGS++))
fi
echo ""

echo "Checking Docker availability..."
echo "-------------------------------"
if command -v docker &> /dev/null; then
    echo -e "${GREEN}✓${NC} Docker is installed"
    docker --version
else
    echo -e "${YELLOW}⚠${NC} Docker not found - install Docker 20.10+"
    ((WARNINGS++))
fi

if command -v docker-compose &> /dev/null; then
    echo -e "${GREEN}✓${NC} Docker Compose is installed"
    docker-compose --version
else
    echo -e "${YELLOW}⚠${NC} Docker Compose not found - install Docker Compose 2.0+"
    ((WARNINGS++))
fi
echo ""

echo "======================================"
echo "Verification Summary"
echo "======================================"
echo ""
echo "Errors: $ERRORS"
echo "Warnings: $WARNINGS"
echo ""

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✓ All checks passed! Docker setup is complete.${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Review and update .env file with secure passwords"
    echo "2. Run: docker-compose build"
    echo "3. Run: docker-compose up -d"
    echo "4. Verify: docker-compose ps"
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo -e "${YELLOW}⚠ Setup is mostly complete but has warnings.${NC}"
    echo "Review warnings above before deployment."
    exit 0
else
    echo -e "${RED}✗ Setup is incomplete. Fix errors above.${NC}"
    exit 1
fi
