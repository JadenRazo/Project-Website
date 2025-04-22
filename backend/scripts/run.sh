#!/bin/bash
set -e

# Print colorful messages
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}Project-Website Server${NC}"
echo -e "${BLUE}========================================${NC}"

# Navigate to project root directory
cd "$(dirname "${BASH_SOURCE[0]}")/../.."
ROOT_DIR="$(pwd)"
BACKEND_DIR="${ROOT_DIR}/backend"
FRONTEND_DIR="${ROOT_DIR}/frontend"

echo -e "${GREEN}Project root directory:${NC} ${YELLOW}${ROOT_DIR}${NC}"
echo -e "${GREEN}Backend directory:${NC} ${YELLOW}${BACKEND_DIR}${NC}"
echo -e "${GREEN}Frontend directory:${NC} ${YELLOW}${FRONTEND_DIR}${NC}"

# Check for required tools and install if missing
echo -e "\n${GREEN}Checking dependencies...${NC}"

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

# Parse command line arguments
ENVIRONMENT="development"
USE_WATCH=false
DEBUG_MODE=false
INSTALL_DEPS=false
SETUP_CONFIG=false

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

# Install dependencies if requested
if [ "$INSTALL_DEPS" = true ]; then
    echo -e "\n${GREEN}Installing backend dependencies...${NC}"
    cd "${BACKEND_DIR}"
    go mod tidy
    go mod download
    
    echo -e "\n${GREEN}Installing frontend dependencies...${NC}"
    cd "${FRONTEND_DIR}"
    npm install
fi

# Setup configuration
CONFIG_FILE="${BACKEND_DIR}/config/${ENVIRONMENT}.yaml"
if [ "$SETUP_CONFIG" = true ] || [ ! -f "$CONFIG_FILE" ]; then
    echo -e "\n${GREEN}Setting up configuration...${NC}"
    
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
  driver: "sqlite"
  dsn: "./database.db"
  maxOpenConns: 25
  maxIdleConns: 25
  connMaxLifetime: 300

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
EOF
    fi
    
    echo -e "${GREEN}Configuration file created at:${NC} ${YELLOW}${CONFIG_FILE}${NC}"
    echo -e "${YELLOW}Please review and update the configuration before deploying to production.${NC}"
fi

# Set the CONFIG_PATH environment variable
export CONFIG_PATH="${CONFIG_FILE}"
echo -e "${GREEN}Using config:${NC} ${YELLOW}${CONFIG_PATH}${NC}"

# Check for database migrations
echo -e "\n${GREEN}Checking database migrations...${NC}"
MIGRATION_CMD="${BACKEND_DIR}/cmd/migration/main.go"
if [ -f "$MIGRATION_CMD" ]; then
    echo -e "${YELLOW}Running database migrations...${NC}"
    go run "$MIGRATION_CMD"
else
    echo -e "${YELLOW}Migration command not found. Skipping migrations.${NC}"
fi

# Check backend main file
MAIN_CMD="${BACKEND_DIR}/cmd/api/main.go"
if [ ! -f "$MAIN_CMD" ]; then
    echo -e "${RED}Error: Main application file not found at ${MAIN_CMD}${NC}"
    exit 1
fi

# Kill existing tmux sessions if they exist
echo -e "\n${GREEN}Cleaning up existing sessions...${NC}"
tmux kill-session -t backend 2>/dev/null || true
tmux kill-session -t frontend 2>/dev/null || true

# Set up debug flags if needed
DEBUG_FLAGS=""
if [ "$DEBUG_MODE" = true ]; then
    DEBUG_FLAGS="GODEBUG=gctrace=1"
    echo -e "${YELLOW}Running in debug mode${NC}"
fi

# Start the backend in a tmux session
echo -e "\n${GREEN}Starting backend in tmux session...${NC}"
if [ "$USE_WATCH" = true ]; then
    # Check if air is installed
    if ! command -v air &> /dev/null; then
        echo -e "${YELLOW}Installing air for hot reloading...${NC}"
        go install github.com/cosmtrek/air@latest
        
        # Check if installation was successful
        if ! command -v air &> /dev/null; then
            echo -e "${RED}Failed to install air. Please install it manually:${NC}"
            echo -e "  ${YELLOW}go install github.com/cosmtrek/air@latest${NC}"
            echo -e "${YELLOW}Running without watch mode instead.${NC}"
            USE_WATCH=false
        fi
    fi
    
    # Check if .air.toml exists
    AIR_CONFIG="${BACKEND_DIR}/.air.toml"
    if [ ! -f "$AIR_CONFIG" ]; then
        echo -e "${YELLOW}Creating air configuration...${NC}"
        
        cat > "$AIR_CONFIG" << EOF
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ./cmd/api"
bin = "tmp/main"
include_ext = ["go", "yaml", "yml"]
exclude_dir = ["tmp", "vendor", "frontend"]
include_dir = []
exclude_file = []
delay = 1000
kill_delay = 500
stop_on_error = true

[color]
main = "magenta"
watcher = "cyan"
build = "yellow"
runner = "green"

[log]
time = true

[misc]
clean_on_exit = true
EOF
    fi
    
    # Start backend with air in a tmux session
    tmux new-session -d -s backend -c "${BACKEND_DIR}" "cd ${BACKEND_DIR} && $DEBUG_FLAGS air; read"
else
    # Start backend in a tmux session
    tmux new-session -d -s backend -c "${BACKEND_DIR}" "cd ${BACKEND_DIR} && $DEBUG_FLAGS go run ${MAIN_CMD}; read"
fi

# Start frontend in a tmux session
echo -e "\n${GREEN}Starting frontend in tmux session...${NC}"
tmux new-session -d -s frontend -c "${FRONTEND_DIR}" "cd ${FRONTEND_DIR} && npm run build && sudo -E npm run start:prod; read"

echo -e "\n${GREEN}All services started successfully!${NC}"
echo -e "\n${BLUE}To attach to the sessions:${NC}"
echo -e "  ${YELLOW}tmux attach-session -t backend${NC}  - View backend logs"
echo -e "  ${YELLOW}tmux attach-session -t frontend${NC} - View frontend logs"
echo -e "\n${BLUE}To detach from a session:${NC} Press ${YELLOW}Ctrl+b${NC} then ${YELLOW}d${NC}"
echo -e "\n${BLUE}To stop all sessions:${NC}"
echo -e "  ${YELLOW}tmux kill-session -t backend${NC}"
echo -e "  ${YELLOW}tmux kill-session -t frontend${NC}"
echo -e "  or run this script again to restart everything" 