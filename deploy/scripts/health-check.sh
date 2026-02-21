#!/bin/bash

# =====================================================
# RDP Platform Health Check Script
# Version: 1.0
# Date: 2026-02-22
# =====================================================

set -e

# Configuration
RDP_HOST="localhost"
API_PORT="8080"
CASDOOR_PORT="8000"
NGINX_PORT="80"
DB_NAME="rdp"
DB_USER="rdp_user"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

ERRORS=0
WARNINGS=0

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_ok() {
    echo -e "${GREEN}[OK]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
    WARNINGS=$((WARNINGS + 1))
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
    ERRORS=$((ERRORS + 1))
}

check_service() {
    local service=$1
    if systemctl is-active --quiet "$service"; then
        log_ok "$service is running"
    else
        log_error "$service is NOT running"
    fi
}

check_port() {
    local host=$1
    local port=$2
    local name=$3
    
    if timeout 2 bash -c "cat < /dev/null > /dev/tcp/$host/$port" 2>/dev/null; then
        log_ok "$name (port $port) is accessible"
    else
        log_error "$name (port $port) is NOT accessible"
    fi
}

check_api() {
    local response
    response=$(curl -s -o /dev/null -w "%{http_code}" "http://$RDP_HOST:$API_PORT/api/v1/health" 2>/dev/null || echo "000")
    
    if [ "$response" = "200" ]; then
        log_ok "API health check passed"
    else
        log_error "API health check failed (HTTP $response)"
    fi
}

check_database() {
    if PGPASSWORD=rdp_password psql -h localhost -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1;" &>/dev/null; then
        log_ok "Database connection successful"
    else
        log_error "Database connection failed"
    fi
}

check_disk() {
    local usage=$(df -h / | awk 'NR==2 {print $5}' | sed 's/%//')
    
    if [ "$usage" -lt 80 ]; then
        log_ok "Disk usage: ${usage}%"
    elif [ "$usage" -lt 90 ]; then
        log_warn "Disk usage: ${usage}% (warning)"
    else
        log_error "Disk usage: ${usage}% (critical)"
    fi
}

check_memory() {
    local available=$(free -m | awk 'NR==2 {print $7}')
    
    if [ "$available" -gt 500 ]; then
        log_ok "Memory available: ${available}MB"
    elif [ "$available" -gt 200 ]; then
        log_warn "Memory available: ${available}MB (low)"
    else
        log_error "Memory available: ${available}MB (critical)"
    fi
}

# =====================================================
# Main
# =====================================================

echo "========================================"
echo "RDP Platform Health Check"
echo "========================================"
echo ""

log_info "Checking system resources..."
check_disk
check_memory
echo ""

log_info "Checking services..."
check_service postgresql
check_service nginx
check_service rdp-api
check_service rdp-casdoor
echo ""

log_info "Checking ports..."
check_port localhost $NGINX_PORT "Nginx"
check_port localhost $API_PORT "API"
check_port localhost $CASDOOR_PORT "Casdoor"
echo ""

log_info "Checking API..."
check_api
echo ""

log_info "Checking database..."
check_database
echo ""

# Summary
echo "========================================"
echo "Summary"
echo "========================================"
echo -e "Errors:   ${RED}$ERRORS${NC}"
echo -e "Warnings: ${YELLOW}$WARNINGS${NC}"

if [ $ERRORS -gt 0 ]; then
    echo ""
    echo "Health check FAILED"
    exit 1
elif [ $WARNINGS -gt 0 ]; then
    echo ""
    echo "Health check passed with warnings"
    exit 0
else
    echo ""
    echo "Health check PASSED"
    exit 0
fi
