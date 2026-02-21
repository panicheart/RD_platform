#!/bin/bash

# =====================================================
# RDP Platform Installation Script
# Version: 1.0
# Date: 2026-02-22
# =====================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
RDP_USER="rdp"
RDP_GROUP="rdp"
RDP_DIR="/opt/rdp"
DB_NAME="rdp"
DB_USER="rdp_user"
DB_PASS="rdp_password"

# =====================================================
# Functions
# =====================================================

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[OK]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_root() {
    if [ "$EUID" -ne 0 ]; then
        log_error "Please run as root"
        exit 1
    fi
}

detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
        VER=$VERSION_ID
    else
        log_error "Cannot detect OS"
        exit 1
    fi
    
    case $OS in
        ubuntu|debian)
            PKG_MGR="apt"
            ;;
        centos|rhel|rocky|alma)
            PKG_MGR="yum"
            ;;
        *)
            log_error "Unsupported OS: $OS"
            exit 1
            ;;
    esac
    
    log_info "Detected OS: $OS $VER"
}

install_dependencies() {
    log_info "Installing system dependencies..."
    
    if [ "$PKG_MGR" = "apt" ]; then
        apt update
        apt install -y curl wget git nginx postgresql postgresql-contrib \
            build-essential libpq-dev libssl-dev
    else
        yum install -y curl wget git nginx postgresql-server postgresql-contrib \
            make gcc libpq-devel openssl-devel
    fi
    
    # Initialize PostgreSQL on RHEL-based systems
    if [ "$PKG_MGR" = "yum" ] && [ ! -d "/var/lib/pgsql/data" ]; then
        postgresql-setup initdb
    fi
    
    log_success "Dependencies installed"
}

create_user() {
    log_info "Creating RDP user..."
    
    if ! id "$RDP_USER" &>/dev/null; then
        useradd -r -s /bin/bash -d "$RDP_DIR" "$RDP_USER" || true
        log_success "User created"
    else
        log_info "User already exists"
    fi
}

create_directories() {
    log_info "Creating directories..."
    
    mkdir -p "$RDP_DIR"/{bin,config,data/{postgresql,files/{projects,knowledge},gitea},logs,web,scripts}
    mkdir -p "$RDP_DIR/data/postgresql"/{data,wal}
    mkdir -p "$RDP_DIR/logs"
    
    chown -R "$RDP_USER:$RDP_GROUP" "$RDP_DIR"
    chmod -R 755 "$RDP_DIR"
    
    log_success "Directories created"
}

configure_postgresql() {
    log_info "Configuring PostgreSQL..."
    
    # Set PostgreSQL password
    su - postgres -c "psql -c \"ALTER USER postgres PASSWORD '$DB_PASS';\"" 2>/dev/null || true
    
    # Create database and user
    su - postgres -c "psql -c \"CREATE USER $DB_USER WITH PASSWORD '$DB_PASS';\"" 2>/dev/null || true
    su - postgres -c "psql -c \"CREATE DATABASE $DB_NAME OWNER $DB_USER;\"" 2>/dev/null || true
    su - postgres -c "psql -c \"GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;\"" 2>/dev/null || true
    
    # Enable UUID extension
    su - postgres -c "psql -d $DB_NAME -c 'CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";'" 2>/dev/null || true
    su - postgres -c "psql -d $DB_NAME -c 'CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";'" 2>/dev/null || true
    
    log_success "PostgreSQL configured"
}

initialize_database() {
    log_info "Initializing database schema..."
    
    if [ -f "$RDP_DIR/database/init.sql" ]; then
        su - postgres -c "psql -d $DB_NAME -f $RDP_DIR/database/init.sql" || log_warn "Some errors occurred during initialization"
        log_success "Database initialized"
    else
        log_error "Database init.sql not found"
    fi
}

install_systemd_services() {
    log_info "Installing systemd services..."
    
    # Copy service files
    cp "$RDP_DIR/deploy/systemd/"*.service /etc/systemd/system/ 2>/dev/null || true
    
    # Reload systemd
    systemctl daemon-reload
    
    # Enable services
    systemctl enable rdp-api 2>/dev/null || true
    systemctl enable rdp-casdoor 2>/dev/null || true
    
    log_success "Systemd services installed"
}

install_nginx_config() {
    log_info "Installing Nginx configuration..."
    
    # Copy nginx config
    cp "$RDP_DIR/deploy/nginx/nginx.conf" /etc/nginx/nginx.conf
    
    # Copy site config
    mkdir -p /etc/nginx/sites-available /etc/nginx/sites-enabled
    cp "$RDP_DIR/deploy/nginx/sites-available/rdp.conf" /etc/nginx/sites-available/
    
    # Enable site
    ln -sf /etc/nginx/sites-available/rdp.conf /etc/nginx/sites-enabled/rdp.conf
    
    # Test nginx config
    nginx -t
    
    log_success "Nginx configured"
}

create_config_files() {
    log_info "Creating configuration files..."
    
    # API config
    cat > "$RDP_DIR/config/rdp-api.yaml" << EOF
app:
  name: "RDP API"
  version: "1.0.0"
  mode: "production"

server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 60
  write_timeout: 60

database:
  host: "localhost"
  port: 5432
  name: "$DB_NAME"
  user: "$DB_USER"
  password: "$DB_PASS"
  ssl_mode: "disable"
  max_open_conns: 25
  max_idle_conns: 5

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

casdoor:
  endpoint: "http://localhost:8000"
  client_id: "rdp-client"
  client_secret: ""
  jwt_secret: "your-jwt-secret-key-change-in-production"

log:
  level: "info"
  format: "json"
  output: "/opt/rdp/logs/api.log"

security:
  session_timeout: 1800
  max_login_attempts: 5
  lockout_duration: 1800
EOF

    # Casdoor config
    cat > "$RDP_DIR/config/casdoor.conf" << EOF
appname = Casdoor
httpport = 8000
version = 1.0.0
logDir = /opt/rdp/casdoor/logs
dataDir = /opt/rdp/casdoor/data

[database]
type = postgres
host = localhost
port = 5432
user = $DB_USER
password = $DB_PASS
database = $DB_NAME

[redis]
endpoint = localhost:6379
password =
db = 0
EOF

    chown -R "$RDP_USER:$RDP_GROUP" "$RDP_DIR/config"
    log_success "Configuration files created"
}

start_services() {
    log_info "Starting services..."
    
    # Start PostgreSQL
    if command -v systemctl &> /dev/null; then
        systemctl start postgresql
        systemctl enable postgresql
    fi
    
    # Wait for PostgreSQL
    sleep 3
    
    # Start Nginx
    systemctl restart nginx
    systemctl enable nginx
    
    # Note: Binary files not yet built, showing instructions
    log_warn "Note: API and Casdoor binaries need to be built first"
    log_info "After building binaries, run:"
    log_info "  systemctl start rdp-api"
    log_info "  systemctl start rdp-casdoor"
}

print_summary() {
    echo ""
    echo "========================================"
    echo -e "${GREEN}RDP Platform Installation Complete${NC}"
    echo "========================================"
    echo ""
    echo "Next steps:"
    echo "  1. Build the API binary: cd services/api && go build -o $RDP_DIR/bin/rdp-api"
    echo "  2. Download Casdoor binary or build from source"
    echo "  3. Build frontend: cd apps/web && npm run build"
    echo "  4. Start services: systemctl start rdp-api rdp-casdoor"
    echo ""
    echo "Access:"
    echo "  Frontend: http://<server-ip>/"
    echo "  API:      http://<server-ip>:8080"
    echo "  Casdoor:  http://<server-ip>:8000"
    echo ""
    echo "Default credentials:"
    echo "  Username: admin"
    echo "  Password: admin123 (MUST CHANGE ON FIRST LOGIN)"
    echo ""
}

# =====================================================
# Main
# =====================================================

main() {
    echo "========================================"
    echo "RDP Platform Installation Script"
    echo "========================================"
    
    check_root
    detect_os
    install_dependencies
    create_user
    create_directories
    configure_postgresql
    initialize_database
    install_systemd_services
    install_nginx_config
    create_config_files
    start_services
    print_summary
}

main "$@"
