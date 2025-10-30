#!/bin/bash
set -e

# Print colorful messages
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Function to print a header
print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${GREEN}$1${NC}"
echo -e "${BLUE}========================================${NC}"
}

# Function to check if a service is running
check_service() {
    local service_name=$1
    local port=$2
    
    echo -e "${YELLOW}Checking $service_name on port $port...${NC}"
    if nc -z localhost $port >/dev/null 2>&1; then
        echo -e "${GREEN}✓ $service_name is running on port $port${NC}"
        return 0
    else
        echo -e "${RED}✗ $service_name is not running on port $port${NC}"
        return 1
    fi
}

# Function to check database connection
check_database() {
    echo -e "${YELLOW}Checking database connection...${NC}"
    
    if [ -f "${BACKEND_DIR}/.env" ]; then
        # Source the env file values for DB connection if available
        source <(grep -E '^DB_' "${BACKEND_DIR}/.env" | sed 's/^/export /')
    fi
    
    # Default to PostgreSQL if no DB specified
    DB_HOST=${DB_HOST:-"localhost"}
    DB_PORT=${DB_PORT:-"5432"}
    DB_USER=${DB_USER:-"postgres"}
    DB_NAME=${DB_NAME:-"project_website"}
    
    # Try to connect to the database
    if command -v psql &> /dev/null; then
        if PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c '\q' &>/dev/null; then
            echo -e "${GREEN}✓ Database connection successful${NC}"
            return 0
        else
            echo -e "${RED}✗ Failed to connect to the database${NC}"
            echo -e "${YELLOW}Please check your database configuration in .env or config files${NC}"
            return 1
        fi
    else
        echo -e "${YELLOW}psql client not found, skipping database connection check${NC}"
        echo -e "${YELLOW}Install with: sudo apt-get install postgresql-client${NC}"
        return 0
    fi
}

# Function to check if Redis is running
check_redis() {
    echo -e "${YELLOW}Checking Redis connection...${NC}"
    
    if command -v redis-cli &> /dev/null; then
        if redis-cli ping &>/dev/null; then
            echo -e "${GREEN}✓ Redis connection successful${NC}"
            return 0
        else
            echo -e "${RED}✗ Failed to connect to Redis${NC}"
            echo -e "${YELLOW}Make sure Redis server is running${NC}"
            return 1
        fi
    else
        echo -e "${YELLOW}redis-cli not found, skipping Redis connection check${NC}"
        echo -e "${YELLOW}Install with: sudo apt-get install redis-tools${NC}"
        return 0
    fi
}

print_header "Project-Website Server"

# Navigate to project root directory
cd "$(dirname "${BASH_SOURCE[0]}")/../.."
ROOT_DIR="$(pwd)"
BACKEND_DIR="${ROOT_DIR}/backend"
FRONTEND_DIR="${ROOT_DIR}/frontend"

echo -e "${GREEN}Project root directory:${NC} ${YELLOW}${ROOT_DIR}${NC}"
echo -e "${GREEN}Backend directory:${NC} ${YELLOW}${BACKEND_DIR}${NC}"
echo -e "${GREEN}Frontend directory:${NC} ${YELLOW}${FRONTEND_DIR}${NC}"

# Check for required tools and install if missing
print_header "Checking dependencies"

# Check for tmux
if ! command -v tmux &> /dev/null; then
    echo -e "${YELLOW}tmux not found. Attempting to install...${NC}"
    if command -v apt-get &> /dev/null; then
        sudo apt-get update && sudo apt-get install -y tmux
    elif command -v yum &> /dev/null; then
        sudo yum install -y tmux
    elif command -v brew &> /dev/null; then
        brew install tmux
    else
        echo -e "${RED}Error: Could not install tmux. Please install it manually.${NC}"
        exit 1
    fi
fi

# Check for Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go first.${NC}"
    echo -e "  Visit ${YELLOW}https://go.dev/doc/install${NC} for installation instructions."
    exit 1
fi

# Check for Node.js and npm
if ! command -v node &> /dev/null || ! command -v npm &> /dev/null; then
    echo -e "${RED}Error: Node.js or npm is not installed. Please install them first.${NC}"
    echo -e "  Visit ${YELLOW}https://nodejs.org/${NC} for installation instructions."
    exit 1
fi

# Check for nc (netcat) for service checking
if ! command -v nc &> /dev/null; then
    echo -e "${YELLOW}netcat not found. Attempting to install...${NC}"
    if command -v apt-get &> /dev/null; then
        sudo apt-get update && sudo apt-get install -y netcat-openbsd
    elif command -v yum &> /dev/null; then
        sudo yum install -y nc
    elif command -v brew &> /dev/null; then
        brew install netcat
    else
        echo -e "${YELLOW}Warning: Could not install netcat. Service status checks will be disabled.${NC}"
    fi
fi

# Parse command line arguments
ENVIRONMENT="development"
USE_WATCH=false
DEBUG_MODE=false
INSTALL_DEPS=false
SETUP_CONFIG=false
PRODUCTION_MODE=false

print_usage() {
    echo -e "${GREEN}Usage:${NC}"
    echo -e "  ${YELLOW}./run.sh [options]${NC}"
    echo -e ""
    echo -e "${GREEN}Options:${NC}"
    echo -e "  ${YELLOW}-e, --env ENV${NC}        Run in specific environment (development, staging, production)"
    echo -e "  ${YELLOW}-w, --watch${NC}          Run backend with file watching (requires air)"
    echo -e "  ${YELLOW}-d, --debug${NC}          Run in debug mode"
    echo -e "  ${YELLOW}-i, --install${NC}        Install/update dependencies"
    echo -e "  ${YELLOW}-s, --setup${NC}          Setup configuration files"
    echo -e "  ${YELLOW}-p, --production${NC}     Run in production mode"
    echo -e "  ${YELLOW}-h, --help${NC}           Print this help message"
}

while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--env)
            ENVIRONMENT="$2"
            shift 2
            ;;
        -w|--watch)
            USE_WATCH=true
            shift
            ;;
        -d|--debug)
            DEBUG_MODE=true
            shift
            ;;
        -i|--install)
            INSTALL_DEPS=true
            shift
            ;;
        -s|--setup)
            SETUP_CONFIG=true
            shift
            ;;
        -p|--production)
            PRODUCTION_MODE=true
            ENVIRONMENT="production"
            shift
            ;;
        -h|--help)
            print_usage
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            print_usage
            exit 1
            ;;
    esac
done

echo -e "${GREEN}Environment:${NC} ${YELLOW}${ENVIRONMENT}${NC}"
if [ "$PRODUCTION_MODE" = true ]; then
    echo -e "${YELLOW}Running in PRODUCTION mode${NC}"
fi

# Install dependencies if requested
if [ "$INSTALL_DEPS" = true ]; then
    print_header "Installing dependencies"
    
    echo -e "${GREEN}Installing backend dependencies...${NC}"
    cd "${BACKEND_DIR}"
    go mod tidy
    go mod download
    
    echo -e "${GREEN}Installing frontend dependencies...${NC}"
    cd "${FRONTEND_DIR}"
    npm install
fi

# Setup configuration
CONFIG_FILE="${BACKEND_DIR}/config/${ENVIRONMENT}.yaml"
if [ "$SETUP_CONFIG" = true ] || [ ! -f "$CONFIG_FILE" ]; then
    print_header "Setting up configuration"
    
    # Check if example configuration exists
    EXAMPLE_CONFIG="${BACKEND_DIR}/config/${ENVIRONMENT}.yaml.example"
    if [ -f "$EXAMPLE_CONFIG" ]; then
        echo -e "${YELLOW}Creating configuration from example...${NC}"
        cp "$EXAMPLE_CONFIG" "$CONFIG_FILE"
    else
        # Create a basic configuration file
        echo -e "${YELLOW}Creating basic configuration file...${NC}"
        mkdir -p "$(dirname "$CONFIG_FILE")"
        
        cat > "$CONFIG_FILE" << EOF
# ${ENVIRONMENT} configuration
server:
  host: "localhost"
  port: 8080
  shutdownTimeout: 30
  readTimeout: 15
  writeTimeout: 15
  idleTimeout: 60

database:
  driver: "postgres"
  dsn: "host=localhost port=5432 user=postgres password=your_secure_password dbname=project_website sslmode=disable"
  maxOpenConns: 25
  maxIdleConns: 25
  connMaxLifetime: 300

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

auth:
  jwtSecret: "change-this-to-a-secure-secret"
  tokenExpiry: 86400
  refreshTokenExpiry: 604800

logging:
  level: "info"
  format: "text"
  filePath: "./logs/backend.log"
  maxSize: 100
  maxBackups: 5
  maxAge: 30
  compress: true

services:
  urlShortener: "http://localhost:8083"
  messaging: "http://localhost:8082"
  devpanel: "http://localhost:8081"
EOF
    fi
    
    echo -e "${GREEN}Configuration file created at:${NC} ${YELLOW}${CONFIG_FILE}${NC}"
    echo -e "${YELLOW}Please review and update the configuration before deploying to production.${NC}"
fi

# Set the CONFIG_PATH environment variable
export CONFIG_PATH="${CONFIG_FILE}"
echo -e "${GREEN}Using config:${NC} ${YELLOW}${CONFIG_PATH}${NC}"

# Check infrastructure services
print_header "Checking infrastructure services"

# Check database connection
check_database

# Check Redis connection
check_redis

# Ensure logs directory exists
mkdir -p "${ROOT_DIR}/logs"

# Check for database schema
print_header "Database setup"
SCHEMA_FILE="${BACKEND_DIR}/schema.sql"
if [ -f "$SCHEMA_FILE" ]; then
    echo -e "${GREEN}Database schema found at: ${SCHEMA_FILE}${NC}"
    echo -e "${YELLOW}To initialize database, run: psql -d project_website -f ${SCHEMA_FILE}${NC}"
else
    echo -e "${YELLOW}No schema.sql found. Database must be initialized manually.${NC}"
fi

# Check backend main files
print_header "Checking service entry points"

# Define all services and their ports
declare -A SERVICES=(
    ["api"]="8080"
    ["devpanel"]="8081"
    ["messaging"]="8082" 
    ["urlshortener"]="8083"
)

# Check each service's main file
MISSING_SERVICES=false
for service in "${!SERVICES[@]}"; do
    SERVICE_CMD="${BACKEND_DIR}/cmd/${service}/main.go"
    if [ -f "$SERVICE_CMD" ]; then
        echo -e "${GREEN}✓ ${service} service found at ${SERVICE_CMD}${NC}"
    else
        echo -e "${YELLOW}✗ ${service} service not found at ${SERVICE_CMD}${NC}"
        MISSING_SERVICES=true
    fi
done

if [ "$MISSING_SERVICES" = true ]; then
    echo -e "${YELLOW}Warning: Some services are missing. The application might not function fully.${NC}"
    echo -e "${YELLOW}Continuing with available services...${NC}"
fi

# Kill existing tmux sessions if they exist
print_header "Cleaning up existing sessions"
for service in "${!SERVICES[@]}"; do
    tmux kill-session -t "$service" 2>/dev/null || true
done
tmux kill-session -t frontend 2>/dev/null || true

# Set up debug flags if needed
DEBUG_FLAGS=""
if [ "$DEBUG_MODE" = true ]; then
    DEBUG_FLAGS="GODEBUG=gctrace=1"
    echo -e "${YELLOW}Running in debug mode${NC}"
fi

# Start the backend services in tmux sessions
print_header "Starting backend services"

# Function to start a service
start_service() {
    local service=$1
    local port=$2
    local cmd_path="${BACKEND_DIR}/cmd/${service}/main.go"
    
    if [ ! -f "$cmd_path" ]; then
        echo -e "${YELLOW}Skipping ${service} service - main.go not found${NC}"
        return
    fi
    
    echo -e "${CYAN}Starting ${service} service on port ${port}...${NC}"
    
if [ "$USE_WATCH" = true ]; then
        # Start with air if watch mode is enabled
        tmux new-session -d -s "$service" -c "${BACKEND_DIR}" \
            "cd ${BACKEND_DIR} && PORT=${port} $DEBUG_FLAGS air -c .air.toml -- ${cmd_path}; read"
    else
        # Start normally
        tmux new-session -d -s "$service" -c "${BACKEND_DIR}" \
            "cd ${BACKEND_DIR} && PORT=${port} $DEBUG_FLAGS go run ${cmd_path}; read"
    fi
    
    # Give the service a moment to start
    sleep 2
    
    # Check if service started successfully
    if check_service "$service" "$port"; then
        echo -e "${GREEN}✓ ${service} service started successfully${NC}"
    else
        echo -e "${YELLOW}⚠ ${service} service might have issues, check logs with:${NC}"
        echo -e "  ${BLUE}tmux attach-session -t ${service}${NC}"
    fi
}

# Start each service
for service in "${!SERVICES[@]}"; do
    start_service "$service" "${SERVICES[$service]}"
done

# Start frontend
print_header "Starting frontend"

# Skip frontend start if requested by sourcing script
if [ "${SKIP_FRONTEND_START}" = true ]; then
    echo -e "${YELLOW}Skipping frontend start as requested by SKIP_FRONTEND_START=true${NC}"
else
    # Determine how to start the frontend based on production mode
    if [ "$PRODUCTION_MODE" = true ]; then
        echo -e "${CYAN}Building frontend for production...${NC}"
        tmux new-session -d -s frontend -c "${FRONTEND_DIR}" \
            "cd ${FRONTEND_DIR} && npm run build:prod && sudo -E npm run start:prod; read"
    else
        echo -e "${CYAN}Starting frontend in development mode...${NC}"
        tmux new-session -d -s frontend -c "${FRONTEND_DIR}" \
            "cd ${FRONTEND_DIR} && npm run start; read"
    fi

    # Give the frontend a moment to start
    sleep 3

    # Check if frontend started
    if check_service "frontend" "3000" || check_service "frontend" "80"; then
        echo -e "${GREEN}✓ Frontend started successfully${NC}"
    else
        echo -e "${YELLOW}⚠ Frontend might have issues, check logs with:${NC}"
        echo -e "  ${BLUE}tmux attach-session -t frontend${NC}"
    fi
fi

# Final status report
print_header "Service Status Report"

# Check each service again and create a summary
all_services_running=true
echo -e "${CYAN}Backend Services:${NC}"
for service in "${!SERVICES[@]}"; do
    port="${SERVICES[$service]}"
    if nc -z localhost "$port" >/dev/null 2>&1; then
        echo -e "  ${GREEN}✓ ${service} (port ${port})${NC}: Running"
    else
        echo -e "  ${RED}✗ ${service} (port ${port})${NC}: Not running"
        all_services_running=false
    fi
done

echo -e "\n${CYAN}Frontend:${NC}"
if nc -z localhost 3000 >/dev/null 2>&1; then
    echo -e "  ${GREEN}✓ Development server (port 3000)${NC}: Running"
elif nc -z localhost 80 >/dev/null 2>&1; then
    echo -e "  ${GREEN}✓ Production server (port 80)${NC}: Running"
else
    echo -e "  ${RED}✗ Frontend server${NC}: Not running"
    all_services_running=false
fi

if [ "$all_services_running" = true ]; then
    print_header "All services started successfully! 🚀"
else
    print_header "Some services might have issues. Check the logs for more details."
fi

# Print helpful information
echo -e "\n${BLUE}Access your application:${NC}"
if [ "$PRODUCTION_MODE" = true ]; then
    echo -e "  ${YELLOW}Frontend:${NC} http://localhost"
else
    echo -e "  ${YELLOW}Frontend:${NC} http://localhost:3000"
fi
echo -e "  ${YELLOW}API:${NC} http://localhost:8080"
echo -e "  ${YELLOW}DevPanel:${NC} http://localhost:8081"
echo -e "  ${YELLOW}Messaging:${NC} http://localhost:8082"
echo -e "  ${YELLOW}URL Shortener:${NC} http://localhost:8083"

echo -e "\n${BLUE}Service management:${NC}"
echo -e "  ${YELLOW}View service logs:${NC}"
for service in "${!SERVICES[@]}"; do
    echo -e "    - ${service}: ${BLUE}tmux attach-session -t ${service}${NC}"
done
echo -e "    - Frontend: ${BLUE}tmux attach-session -t frontend${NC}"
echo -e "  ${YELLOW}Detach from tmux:${NC} Press ${BLUE}Ctrl+b${NC} then ${BLUE}d${NC}"
echo -e "  ${YELLOW}Stop all services:${NC} Run this script with ${BLUE}./stop.sh${NC} (if available) or kill tmux sessions"

echo -e "\n${MAGENTA}Happy coding! 🚀${NC}" 