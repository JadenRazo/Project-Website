#!/bin/bash

# Print colorful messages
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print a header
print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${GREEN}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_header "Project-Website Service Status"

# Define all services
SERVICES=("api" "devpanel" "messaging" "urlshortener" "frontend")

# Services ports
declare -A SERVICE_PORTS
SERVICE_PORTS=([api]=8000 [devpanel]=8090 [messaging]=8070 [urlshortener]=8080 [frontend]=3000)

# Service descriptions
declare -A SERVICE_DESCRIPTIONS
SERVICE_DESCRIPTIONS=(
    [api]="Main API Service"
    [devpanel]="Developer Panel"
    [messaging]="Messaging Service"
    [urlshortener]="URL Shortener"
    [frontend]="Frontend Application"
)

# Function to check if a service is running in tmux
check_tmux_service() {
    local service=$1
    if tmux has-session -t "$service" 2>/dev/null; then
        return 0
    else
        return 1
    fi
}

# Function to check if a service port is open
check_service_port() {
    local port=$1
    if command -v nc &> /dev/null; then
        nc -z localhost "$port" &>/dev/null
        return $?
    elif command -v curl &> /dev/null; then
        curl -s "http://localhost:$port" -o /dev/null
        local status=$?
        if [ $status -eq 7 ]; then  # Connection refused
            return 1
        elif [ $status -eq 52 ] || [ $status -eq 56 ]; then  # Empty response or recv failure
            return 0  # Port is open but service might not be fully running
        else
            return 0  # Any other result means port is open
        fi
    else
        # Fallback if neither nc nor curl are available
        (echo > /dev/tcp/localhost/"$port") &>/dev/null
        return $?
    fi
}

# Function to check database
check_database() {
    local status="Unknown"
    
    if command -v psql &> /dev/null; then
        echo -e "${CYAN}Checking PostgreSQL database...${NC}"
        
        if [ -f "../db_config.yml" ]; then
            DB_USER=$(grep -E "^user:" ../db_config.yml | cut -d ":" -f2 | tr -d " ")
            DB_PASS=$(grep -E "^password:" ../db_config.yml | cut -d ":" -f2 | tr -d " ")
            DB_NAME=$(grep -E "^dbname:" ../db_config.yml | cut -d ":" -f2 | tr -d " ")
            DB_HOST=$(grep -E "^host:" ../db_config.yml | cut -d ":" -f2 | tr -d " ")
            DB_PORT=$(grep -E "^port:" ../db_config.yml | cut -d ":" -f2 | tr -d " ")
            
            if [ -z "$DB_USER" ] || [ -z "$DB_NAME" ] || [ -z "$DB_HOST" ]; then
                echo -e "${YELLOW}Incomplete database configuration. Using defaults.${NC}"
                PGPASSWORD="postgres" psql -h localhost -U postgres -c "SELECT 1;" &>/dev/null
            else
                PGPASSWORD="$DB_PASS" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1;" &>/dev/null
            fi
            
            if [ $? -eq 0 ]; then
                status="${GREEN}Connected${NC}"
            else
                status="${RED}Not Connected${NC}"
            fi
        else
            # Try default connection parameters
            PGPASSWORD="postgres" psql -h localhost -U postgres -c "SELECT 1;" &>/dev/null
            if [ $? -eq 0 ]; then
                status="${GREEN}Connected${NC}"
            else
                status="${RED}Not Connected${NC}"
            fi
        fi
    else
        status="${YELLOW}psql not installed${NC}"
    fi
    
    echo -e "PostgreSQL: $status"
}

# Function to check Redis
check_redis() {
    local status="Unknown"
    
    if command -v redis-cli &> /dev/null; then
        echo -e "${CYAN}Checking Redis...${NC}"
        
        if redis-cli ping &>/dev/null; then
            status="${GREEN}Connected${NC}"
        else
            status="${RED}Not Connected${NC}"
        fi
    else
        status="${YELLOW}redis-cli not installed${NC}"
    fi
    
    echo -e "Redis: $status"
}

# Check if tmux is installed
if ! command -v tmux &> /dev/null; then
    echo -e "${RED}Error: tmux is not installed. Cannot check service status.${NC}"
    exit 1
fi

# Print service status
echo -e "${CYAN}Service Status:${NC}"
echo -e "-------------------------------------"
printf "%-15s %-15s %-15s %-25s\n" "SERVICE" "TMUX SESSION" "PORT STATUS" "DESCRIPTION"
echo -e "-------------------------------------"

for service in "${SERVICES[@]}"; do
    # Check tmux session
    if check_tmux_service "$service"; then
        tmux_status="${GREEN}Running${NC}"
    else
        tmux_status="${RED}Stopped${NC}"
    fi
    
    # Check port status
    if [ -n "${SERVICE_PORTS[$service]}" ]; then
        if check_service_port "${SERVICE_PORTS[$service]}"; then
            port_status="${GREEN}Open${NC}"
        else
            port_status="${RED}Closed${NC}"
        fi
    else
        port_status="${YELLOW}Unknown${NC}"
    fi
    
    # Print status line
    printf "%-15s ${tmux_status} ${port_status} %-25s\n" "$service" "${SERVICE_DESCRIPTIONS[$service]}"
done

echo -e "\n${CYAN}Infrastructure Services:${NC}"
check_database
check_redis

# Print project URLs
print_header "Project URLs"
echo -e "${YELLOW}Frontend:${NC} http://localhost:3000"
echo -e "${YELLOW}API:${NC} http://localhost:8000"
echo -e "${YELLOW}DevPanel:${NC} http://localhost:8090"
echo -e "${YELLOW}Messaging:${NC} http://localhost:8070"
echo -e "${YELLOW}URL Shortener:${NC} http://localhost:8080" 