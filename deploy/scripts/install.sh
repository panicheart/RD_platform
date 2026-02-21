#!/bin/bash
set -e

RDP_DIR="/opt/rdp"
RDP_USER="rdp"
RDP_GROUP="rdp"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root"
        exit 1
    fi
}

install_dependencies() {
    log_info "Installing dependencies..."
    
    if command -v apt-get &> /dev/null; then
        apt-get update
        apt-get install -y curl wget git nginx postgresql postgresql-contrib
    elif command -v yum &> /dev/null; then
        yum install -y curl wget git nginx postgresql postgresql-server
        postgresql-setup initdb
    else
        log_error "Unsupported OS. Please install dependencies manually."
        exit 1
    fi
}

create_user() {
    log_info "Creating rdp user..."
    if ! id -u $RDP_USER &>/dev/null; then
        useradd -r -s /bin/false $RDP_USER
    fi
}

setup_directories() {
    log_info "Setting up directories..."
    mkdir -p $RDP_DIR/{bin,config,logs,data,backup}
    mkdir -p $RDP_DIR/apps/web/dist
    
    chown -R $RDP_USER:$RDP_GROUP $RDP_DIR
}

setup_postgresql() {
    log_info "Setting up PostgreSQL..."
    
    if command -v systemctl &> /dev/null; then
        systemctl enable postgresql
        systemctl start postgresql
    fi
    
    sudo -u postgres psql -c "CREATE DATABASE rdp;" 2>/dev/null || true
    sudo -u postgres psql -c "CREATE USER rdp WITH PASSWORD 'rdp_password';" 2>/dev/null || true
    sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE rdp TO rdp;" 2>/dev/null || true
}

install_systemd_services() {
    log_info "Installing systemd services..."
    
    cp deploy/systemd/*.service /etc/systemd/system/
    systemctl daemon-reload
    systemctl enable rdp-api
    systemctl enable rdp-casdoor
}

setup_nginx() {
    log_info "Setting up Nginx..."
    
    cp deploy/nginx/nginx.conf /etc/nginx/nginx.conf
    cp deploy/nginx/sites-available/rdp.conf /etc/nginx/sites-available/
    ln -sf /etc/nginx/sites-available/rdp.conf /etc/nginx/sites-enabled/rdp.conf
    rm -f /etc/nginx/sites-enabled/default
    
    nginx -t && systemctl reload nginx
}

main() {
    log_info "Starting RDP Platform installation..."
    
    check_root
    install_dependencies
    create_user
    setup_directories
    setup_postgresql
    install_systemd_services
    setup_nginx
    
    log_info "Installation completed!"
    log_info "Please copy your application files to $RDP_DIR"
    log_info "Then run: systemctl start rdp-api rdp-casdoor"
}

main "$@"
