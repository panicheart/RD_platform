#!/bin/bash
# =============================================================================
# RDP Platform Installation Script
# Description: One-click installation for RDP platform
# Version: 1.0
# Created: 2026-02-21
# =============================================================================

set -euo pipefail

# =============================================================================
# Configuration
# =============================================================================
RDP_DIR="/opt/rdp"
RDP_USER="rdp"
RDP_GROUP="rdp"
RDP_SERVICE_DIR="${RDP_DIR}/services"
RDP_APPS_DIR="${RDP_DIR}/apps"
PROJECT_ROOT=""
GO_VERSION="1.22"
NODE_VERSION="20"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# =============================================================================
# Logging Functions
# =============================================================================
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# =============================================================================
# Error Handling
# =============================================================================
handle_error() {
    log_error "Installation failed at step: $1"
    log_error "Error on line $2"
    exit 1
}

trap 'handle_error "$current_step" "$LINENO"' ERR

# =============================================================================
# Pre-flight Checks
# =============================================================================
current_step="Pre-flight Checks"

check_root() {
    log_step "Checking root privileges..."
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root"
        exit 1
    fi
}

detect_os() {
    log_step "Detecting operating system..."
    if [[ -f /etc/os-release ]]; then
        source /etc/os-release
        OS_ID=$ID
        OS_VERSION=$VERSION_ID
        log_info "Detected OS: $OS_ID $OS_VERSION"
    else
        log_error "Unable to detect OS"
        exit 1
    fi
}

check_project_root() {
    log_step "Detecting project root..."
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
    
    if [[ ! -f "${PROJECT_ROOT}/services/api/main.go" ]]; then
        log_error "Cannot find project root. Expected: ${PROJECT_ROOT}"
        exit 1
    fi
    
    log_info "Project root: ${PROJECT_ROOT}"
}

# =============================================================================
# Dependencies Installation
# =============================================================================
current_step="Install Dependencies"

install_dependencies_debian() {
    log_step "Installing dependencies (Debian/Ubuntu)..."
    
    apt-get update
    apt-get install -y \
        curl \
        wget \
        git \
        nginx \
        postgresql-16 \
        postgresql-client-16 \
        postgresql-contrib-16 \
        redis-server \
        build-essential \
        jq \
        unzip
    
    # Enable and start services
    systemctl enable postgresql
    systemctl enable redis-server
    systemctl enable nginx
    
    systemctl start postgresql
    systemctl start redis-server
    systemctl start nginx
    
    log_info "Dependencies installed successfully"
}

install_dependencies_rhel() {
    log_step "Installing dependencies (RHEL/CentOS)..."
    
    # Add PostgreSQL repository
    dnf install -y https://download.postgresql.org/pub/repos/yum/reporpms/EL-9-x86_64/pgdg-redhat-repo-latest.noarch.rpm
    dnf -qy module disable postgresql
    
    dnf install -y \
        curl \
        wget \
        git \
        nginx \
        postgresql16-server \
        postgresql16-contrib \
        redis \
        gcc \
        make \
        jq \
        unzip
    
    # Initialize PostgreSQL
    /usr/pgsql-16/bin/postgresql-16-setup initdb
    
    # Enable and start services
    systemctl enable postgresql-16
    systemctl enable redis
    systemctl enable nginx
    
    systemctl start postgresql-16
    systemctl start redis
    systemctl start nginx
    
    log_info "Dependencies installed successfully"
}

install_go() {
    log_step "Installing Go ${GO_VERSION}..."
    
    if command -v go &> /dev/null; then
        INSTALLED_GO=$(go version | awk '{print $3}')
        log_info "Go already installed: ${INSTALLED_GO}"
        return 0
    fi
    
    # Download and install Go
    GO_TARBALL="go${GO_VERSION}.linux-amd64.tar.gz"
    cd /tmp
    wget -q "https://go.dev/dl/${GO_TARBALL}"
    rm -rf /usr/local/go
    tar -C /usr/local -xzf "${GO_TARBALL}"
    rm "${GO_TARBALL}"
    
    # Add to PATH
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/go.sh
    
    log_info "Go ${GO_VERSION} installed successfully"
}

install_nodejs() {
    log_step "Installing Node.js ${NODE_VERSION}..."
    
    if command -v node &> /dev/null; then
        INSTALLED_NODE=$(node --version)
        log_info "Node.js already installed: ${INSTALLED_NODE}"
        return 0
    fi
    
    # Install NodeSource repository
    curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash -
    
    if [[ "$OS_ID" == "ubuntu" || "$OS_ID" == "debian" ]]; then
        apt-get install -y nodejs
    else
        dnf install -y nodejs
    fi
    
    log_info "Node.js ${NODE_VERSION} installed successfully"
}

install_dependencies() {
    case "$OS_ID" in
        ubuntu|debian)
            install_dependencies_debian
            ;;
        centos|rhel|rocky|almalinux)
            install_dependencies_rhel
            ;;
        *)
            log_error "Unsupported OS: $OS_ID"
            exit 1
            ;;
    esac
    
    install_go
    install_nodejs
}

# =============================================================================
# User and Directory Setup
# =============================================================================
current_step="Setup User and Directories"

create_user() {
    log_step "Creating rdp user..."
    
    if ! id -u $RDP_USER &>/dev/null; then
        useradd -r -s /bin/false -d $RDP_DIR -m $RDP_USER
        log_info "User $RDP_USER created"
    else
        log_info "User $RDP_USER already exists"
    fi
}

setup_directories() {
    log_step "Setting up directories..."
    
    mkdir -p $RDP_DIR/{bin,config,logs,data,backup,scripts}
    mkdir -p $RDP_DIR/data/{uploads,cache,temp}
    mkdir -p $RDP_DIR/apps/web/dist
    mkdir -p /var/log/rdp
    
    chown -R $RDP_USER:$RDP_GROUP $RDP_DIR
    chown -R $RDP_USER:$RDP_GROUP /var/log/rdp
    
    log_info "Directories created"
}

# =============================================================================
# Database Setup
# =============================================================================
current_step="Setup Database"

setup_postgresql() {
    log_step "Setting up PostgreSQL..."
    
    local db_script="${PROJECT_ROOT}/deploy/scripts/init-db.sh"
    
    if [[ -f "$db_script" ]]; then
        chmod +x "$db_script"
        "$db_script" --install
    else
        log_warn "Database init script not found, using default setup"
        
        sudo -u postgres psql -c "CREATE DATABASE rdp;" 2>/dev/null || true
        sudo -u postgres psql -c "CREATE USER rdp_user WITH PASSWORD 'rdp_secret_2026';" 2>/dev/null || true
        sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE rdp TO rdp_user;" 2>/dev/null || true
    fi
    
    log_info "PostgreSQL setup complete"
}

# =============================================================================
# Build Steps
# =============================================================================
current_step="Build Application"

build_backend() {
    log_step "Building backend API..."
    
    local api_dir="${PROJECT_ROOT}/services/api"
    
    if [[ ! -d "$api_dir" ]]; then
        log_error "API directory not found: $api_dir"
        exit 1
    fi
    
    cd "$api_dir"
    
    # Download dependencies
    log_info "Downloading Go dependencies..."
    export GO111MODULE=on
    export GOPROXY=https://proxy.golang.org,direct
    /usr/local/go/bin/go mod download
    
    # Build binary
    log_info "Compiling API binary..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        /usr/local/go/bin/go build -ldflags="-w -s -X main.Version=$(git describe --tags --always 2>/dev/null || echo 'dev')" \
        -o "${RDP_DIR}/bin/api" main.go
    
    chmod +x "${RDP_DIR}/bin/api"
    chown $RDP_USER:$RDP_GROUP "${RDP_DIR}/bin/api"
    
    log_info "Backend built successfully"
}

build_frontend() {
    log_step "Building frontend..."
    
    local web_dir="${PROJECT_ROOT}/apps/web"
    
    if [[ ! -d "$web_dir" ]]; then
        log_warn "Web directory not found: $web_dir"
        return 0
    fi
    
    cd "$web_dir"
    
    # Install dependencies
    log_info "Installing npm dependencies..."
    npm ci --production=false
    
    # Build
    log_info "Building frontend..."
    npm run build
    
    # Copy dist to deployment directory
    rm -rf "${RDP_DIR}/apps/web/dist"
    cp -r dist "${RDP_DIR}/apps/web/"
    chown -R $RDP_USER:$RDP_GROUP "${RDP_DIR}/apps/web/dist"
    
    log_info "Frontend built successfully"
}

build_all() {
    build_backend
    build_frontend
}

# =============================================================================
# Configuration
# =============================================================================
current_step="Configure Services"

setup_environment() {
    log_step "Setting up environment configuration..."
    
    # Copy environment file if not exists
    if [[ ! -f "${RDP_DIR}/config/rdp-api.env" ]]; then
        if [[ -f "${PROJECT_ROOT}/services/api/.env.production" ]]; then
            cp "${PROJECT_ROOT}/services/api/.env.production" "${RDP_DIR}/config/rdp-api.env"
        else
            cat > "${RDP_DIR}/config/rdp-api.env" << 'EOF'
# RDP API Production Environment
RDP_ENV=production
RDP_SERVER_HOST=0.0.0.0
RDP_API_PORT=8080
RDP_READ_TIMEOUT=30s
RDP_WRITE_TIMEOUT=30s

# Database
RDP_DB_HOST=localhost
RDP_DB_PORT=5432
RDP_DB_USER=rdp_user
RDP_DB_PASSWORD=rdp_secret_2026
RDP_DB_NAME=rdp
RDP_DB_SSLMODE=disable

# JWT (CHANGE THESE IN PRODUCTION!)
RDP_JWT_SECRET=$(openssl rand -base64 32)
RDP_ACCESS_TOKEN_TTL=2h
RDP_REFRESH_TOKEN_TTL=168h
RDP_JWT_ISSUER=rdp-api
RDP_JWT_AUDIENCE=rdp-users

# Log
RDP_LOG_LEVEL=info
RDP_LOG_FORMAT=json
RDP_LOG_OUTPUT=stdout

# Redis
RDP_REDIS_HOST=localhost
RDP_REDIS_PORT=6379
RDP_REDIS_PASSWORD=
RDP_REDIS_DB=0
EOF
        fi
        chown $RDP_USER:$RDP_GROUP "${RDP_DIR}/config/rdp-api.env"
        chmod 600 "${RDP_DIR}/config/rdp-api.env"
        log_info "Environment file created"
    else
        log_info "Environment file already exists"
    fi
}

setup_nginx() {
    log_step "Configuring Nginx..."
    
    # Copy nginx configuration
    cp "${PROJECT_ROOT}/deploy/nginx/nginx.conf" /etc/nginx/nginx.conf
    cp "${PROJECT_ROOT}/deploy/nginx/sites-available/rdp.conf" /etc/nginx/sites-available/
    
    # Enable site
    ln -sf /etc/nginx/sites-available/rdp.conf /etc/nginx/sites-enabled/rdp.conf
    rm -f /etc/nginx/sites-enabled/default
    
    # Test and reload
    if nginx -t; then
        systemctl reload nginx
        log_info "Nginx configured successfully"
    else
        log_error "Nginx configuration test failed"
        exit 1
    fi
}

install_systemd_services() {
    log_step "Installing systemd services..."
    
    cp "${PROJECT_ROOT}/deploy/systemd/"*.service /etc/systemd/system/
    cp "${PROJECT_ROOT}/deploy/systemd/rdp.target" /etc/systemd/system/
    
    # Copy scripts
    cp "${PROJECT_ROOT}/deploy/scripts/health-check.sh" "${RDP_DIR}/scripts/"
    cp "${PROJECT_ROOT}/deploy/scripts/backup.sh" "${RDP_DIR}/scripts/"
    chmod +x "${RDP_DIR}/scripts/"*.sh
    
    systemctl daemon-reload
    systemctl enable rdp.target
    systemctl enable rdp-api
    systemctl enable rdp-web
    
    log_info "Systemd services installed"
}

# =============================================================================
# Firewall Setup
# =============================================================================
current_step="Configure Firewall"

setup_firewall() {
    log_step "Configuring firewall..."
    
    if command -v ufw &> /dev/null; then
        ufw allow 'Nginx Full'
        ufw allow 22/tcp
        log_info "UFW rules added"
    elif command -v firewall-cmd &> /dev/null; then
        firewall-cmd --permanent --add-service=http
        firewall-cmd --permanent --add-service=https
        firewall-cmd --reload
        log_info "Firewalld rules added"
    else
        log_warn "No firewall detected, skipping configuration"
    fi
}

# =============================================================================
# Post-installation
# =============================================================================
current_step="Post-installation"

post_install() {
    log_step "Running post-installation tasks..."
    
    # Create version file
    echo "$(date '+%Y-%m-%d %H:%M:%S')" > "${RDP_DIR}/.install_date"
    echo "1.0.0" > "${RDP_DIR}/.version"
    chown $RDP_USER:$RDP_GROUP "${RDP_DIR}/.install_date" "${RDP_DIR}/.version"
    
    log_info "Post-installation complete"
}

print_summary() {
    echo -e "${GREEN}"
    echo "╔════════════════════════════════════════════════════════════╗"
    echo "║     RDP Platform Installation Complete!                    ║"
    echo "╠════════════════════════════════════════════════════════════╣"
    echo "║  Installation Directory: ${RDP_DIR}"
    echo "║  Web Interface: http://localhost/                          ║"
    echo "║  API Endpoint: http://localhost:8080/                      ║"
    echo "╠════════════════════════════════════════════════════════════╣"
    echo "║  Next Steps:                                               ║"
    echo "║  1. Review config: ${RDP_DIR}/config/rdp-api.env"
    echo "║  2. Start services: systemctl start rdp-api                ║"
    echo "║  3. Check status: systemctl status rdp-api                 ║"
    echo "║  4. View logs: journalctl -u rdp-api -f                    ║"
    echo "╠════════════════════════════════════════════════════════════╣"
    echo "║  Useful Commands:                                          ║"
    echo "║  • systemctl start|stop|restart rdp-api                    ║"
    echo "║  • ${RDP_DIR}/scripts/health-check.sh                      ║"
    echo "║  • ${RDP_DIR}/scripts/backup.sh                            ║"
    echo "╚════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# =============================================================================
# Main
# =============================================================================
main() {
    echo -e "${BLUE}"
    echo "╔════════════════════════════════════════════════════════════╗"
    echo "║     RDP Platform Installation Script                       ║"
    echo "║     Version: 1.0                                           ║"
    echo "╚════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    check_root
    detect_os
    check_project_root
    install_dependencies
    create_user
    setup_directories
    setup_postgresql
    build_all
    setup_environment
    setup_nginx
    install_systemd_services
    setup_firewall
    post_install
    print_summary
}

# Parse arguments
case "${1:-}" in
    --build-only)
        check_project_root
        build_all
        ;;
    --help|-h)
        echo "Usage: $0 [options]"
        echo ""
        echo "Options:"
        echo "  --build-only    Only build the application"
        echo "  --help, -h      Show this help message"
        echo ""
        echo "Environment variables:"
        echo "  RDP_DB_PASSWORD    Set database password"
        exit 0
        ;;
    *)
        main
        ;;
esac
