#!/bin/bash
# =============================================================================
# RDP Database Initialization Script
# Description: Initialize PostgreSQL database for RDP platform
# Version: 1.0
# Created: 2026-02-21
# =============================================================================

set -euo pipefail

# =============================================================================
# Configuration
# =============================================================================
DB_NAME="rdp"
DB_USER="rdp_user"
DB_PASSWORD="${RDP_DB_PASSWORD:-rdp_secret_2026}"
DB_HOST="${RDP_DB_HOST:-localhost}"
DB_PORT="${RDP_DB_PORT:-5432}"
POSTGRES_USER="${POSTGRES_USER:-postgres}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Log file
LOG_DIR="${RDP_LOG_DIR:-/var/log/rdp}"
LOG_FILE="${LOG_DIR}/db-init-$(date +%Y%m%d-%H%M%S).log"

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

# =============================================================================
# Utility Functions
# =============================================================================
log() {
    local level="$1"
    local message="$2"
    local timestamp
    timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    # Write to log file
    mkdir -p "${LOG_DIR}"
    echo "[${timestamp}] [${level}] ${message}" >> "${LOG_FILE}"
    
    # Output to console with colors
    case "${level}" in
        INFO)
            echo -e "${GREEN}[INFO]${NC} ${message}"
            ;;
        WARN)
            echo -e "${YELLOW}[WARN]${NC} ${message}"
            ;;
        ERROR)
            echo -e "${RED}[ERROR]${NC} ${message}"
            ;;
        DEBUG)
            echo -e "${BLUE}[DEBUG]${NC} ${message}"
            ;;
        *)
            echo "[${level}] ${message}"
            ;;
    esac
}

log_info() { log "INFO" "$1"; }
log_warn() { log "WARN" "$1"; }
log_error() { log "ERROR" "$1"; }
log_debug() { log "DEBUG" "$1"; }

# =============================================================================
# Check Functions
# =============================================================================

check_postgres_installed() {
    log_info "检查 PostgreSQL 安装状态..."
    
    if command -v psql &> /dev/null; then
        local version
        version=$(psql --version | head -n1 | awk '{print $3}')
        log_info "PostgreSQL 已安装，版本: ${version}"
        return 0
    fi
    
    # Check if service exists
    if systemctl list-units --type=service | grep -q postgresql; then
        log_warn "PostgreSQL 服务存在但 psql 命令不可用"
        return 1
    fi
    
    log_error "PostgreSQL 未安装"
    return 1
}

check_postgres_running() {
    log_info "检查 PostgreSQL 服务状态..."
    
    if systemctl is-active --quiet postgresql 2>/dev/null || \
       systemctl is-active --quiet postgresql-16 2>/dev/null || \
       pg_isready -h "${DB_HOST}" -p "${DB_PORT}" &>/dev/null; then
        log_info "PostgreSQL 服务正在运行"
        return 0
    fi
    
    log_error "PostgreSQL 服务未运行"
    return 1
}

detect_postgres_service() {
    local services=("postgresql" "postgresql-16" "postgresql@16-main" "postgresql@15-main" "postgresql@14-main")
    
    for service in "${services[@]}"; do
        if systemctl list-units --type=service | grep -q "${service}"; then
            echo "${service}"
            return 0
        fi
    done
    
    echo "postgresql"
}

# =============================================================================
# Installation Functions
# =============================================================================

install_postgres_ubuntu() {
    log_info "在 Ubuntu/Debian 上安装 PostgreSQL..."
    
    apt-get update
    apt-get install -y postgresql-16 postgresql-client-16 postgresql-contrib-16
    
    # Start service
    systemctl enable postgresql
    systemctl start postgresql
    
    log_info "PostgreSQL 16 安装完成"
}

install_postgres_centos() {
    log_info "在 CentOS/RHEL 上安装 PostgreSQL..."
    
    # Add PostgreSQL repository
    if ! rpm -qa | grep -q pgdg; then
        dnf install -y https://download.postgresql.org/pub/repos/yum/reporpms/EL-9-x86_64/pgdg-redhat-repo-latest.noarch.rpm
    fi
    
    # Disable built-in PostgreSQL module
    dnf -qy module disable postgresql
    
    # Install PostgreSQL 16
    dnf install -y postgresql16-server postgresql16-contrib
    
    # Initialize database
    /usr/pgsql-16/bin/postgresql-16-setup initdb
    
    # Start service
    systemctl enable postgresql-16
    systemctl start postgresql-16
    
    log_info "PostgreSQL 16 安装完成"
}

install_postgresql() {
    log_info "开始安装 PostgreSQL..."
    
    if [[ -f /etc/os-release ]]; then
        source /etc/os-release
        case "$ID" in
            ubuntu|debian)
                install_postgres_ubuntu
                ;;
            centos|rhel|rocky|almalinux)
                install_postgres_centos
                ;;
            *)
                log_error "不支持的操作系统: $ID"
                exit 1
                ;;
        esac
    else
        log_error "无法检测操作系统类型"
        exit 1
    fi
}

# =============================================================================
# Database Functions
# =============================================================================

get_psql_cmd() {
    local user="$1"
    local db="${2:-postgres}"
    
    if [[ -n "${POSTGRES_PASSWORD}" ]]; then
        PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${user}" -d "${db}"
    else
        # Try sudo for local postgres user
        if [[ "$user" == "postgres" && "${DB_HOST}" == "localhost" ]]; then
            sudo -u postgres psql -d "${db}"
        else
            psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${user}" -d "${db}"
        fi
    fi
}

create_database() {
    log_info "创建数据库: ${DB_NAME}..."
    
    # Check if database exists
    local db_exists
    if [[ -n "${POSTGRES_PASSWORD}" ]]; then
        db_exists=$(PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${POSTGRES_USER}" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}'" 2>/dev/null || echo "")
    else
        db_exists=$(sudo -u postgres psql -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}'" 2>/dev/null || echo "")
    fi
    
    if [[ "$db_exists" == "1" ]]; then
        log_warn "数据库 ${DB_NAME} 已存在，跳过创建"
        return 0
    fi
    
    # Create database
    if [[ -n "${POSTGRES_PASSWORD}" ]]; then
        PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${POSTGRES_USER}" -d postgres -c "CREATE DATABASE ${DB_NAME} ENCODING 'UTF8' LC_COLLATE 'en_US.UTF-8' LC_CTYPE 'en_US.UTF-8' TEMPLATE template0;" 2>/dev/null
    else
        sudo -u postgres psql -d postgres -c "CREATE DATABASE ${DB_NAME} ENCODING 'UTF8' LC_COLLATE 'en_US.UTF-8' LC_CTYPE 'en_US.UTF-8' TEMPLATE template0;"
    fi
    
    log_info "数据库 ${DB_NAME} 创建成功"
}

create_user() {
    log_info "创建数据库用户: ${DB_USER}..."
    
    # Check if user exists
    local user_exists
    if [[ -n "${POSTGRES_PASSWORD}" ]]; then
        user_exists=$(PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${POSTGRES_USER}" -d postgres -tAc "SELECT 1 FROM pg_roles WHERE rolname='${DB_USER}'" 2>/dev/null || echo "")
    else
        user_exists=$(sudo -u postgres psql -d postgres -tAc "SELECT 1 FROM pg_roles WHERE rolname='${DB_USER}'" 2>/dev/null || echo "")
    fi
    
    if [[ "$user_exists" == "1" ]]; then
        log_warn "用户 ${DB_USER} 已存在，更新密码..."
        if [[ -n "${POSTGRES_PASSWORD}" ]]; then
            PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${POSTGRES_USER}" -d postgres -c "ALTER USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}';" 2>/dev/null
        else
            sudo -u postgres psql -d postgres -c "ALTER USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}';"
        fi
    else
        # Create user
        if [[ -n "${POSTGRES_PASSWORD}" ]]; then
            PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${POSTGRES_USER}" -d postgres -c "CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}';" 2>/dev/null
        else
            sudo -u postgres psql -d postgres -c "CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}';"
        fi
        log_info "用户 ${DB_USER} 创建成功"
    fi
    
    # Grant privileges
    log_info "授予用户权限..."
    if [[ -n "${POSTGRES_PASSWORD}" ]]; then
        PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${POSTGRES_USER}" -d "${DB_NAME}" -c "GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};" 2>/dev/null
        PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${POSTGRES_USER}" -d "${DB_NAME}" -c "ALTER USER ${DB_USER} WITH SUPERUSER;" 2>/dev/null || true
    else
        sudo -u postgres psql -d "${DB_NAME}" -c "GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};"
        sudo -u postgres psql -d "${DB_NAME}" -c "ALTER USER ${DB_USER} WITH SUPERUSER;" || true
    fi
}

create_extensions() {
    log_info "创建 PostgreSQL 扩展..."
    
    local extensions_sql="
CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";
CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";
CREATE EXTENSION IF NOT EXISTS \"pg_trgm\";
"
    
    if [[ -n "${POSTGRES_PASSWORD}" ]]; then
        PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${POSTGRES_USER}" -d "${DB_NAME}" -c "${extensions_sql}" 2>/dev/null
    else
        echo "${extensions_sql}" | sudo -u postgres psql -d "${DB_NAME}"
    fi
    
    log_info "扩展创建完成"
}

# =============================================================================
# Migration Functions
# =============================================================================

run_migrations() {
    log_info "执行数据库迁移..."
    
    local migrations_dir="${PROJECT_ROOT}/database/migrations"
    
    if [[ ! -d "${migrations_dir}" ]]; then
        log_warn "迁移目录不存在: ${migrations_dir}"
        return 0
    fi
    
    # Get list of migration files sorted
    local migrations
    migrations=$(ls -1 "${migrations_dir}"/*.sql 2>/dev/null | sort)
    
    if [[ -z "$migrations" ]]; then
        log_warn "没有找到迁移文件"
        return 0
    fi
    
    for migration in $migrations; do
        local filename
        filename=$(basename "$migration")
        log_info "执行迁移: ${filename}"
        
        if [[ -n "${POSTGRES_PASSWORD}" ]]; then
            PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${POSTGRES_USER}" -d "${DB_NAME}" -f "$migration" 2>/dev/null || {
                log_error "迁移失败: ${filename}"
                return 1
            }
        else
            sudo -u postgres psql -d "${DB_NAME}" -f "$migration" || {
                log_error "迁移失败: ${filename}"
                return 1
            }
        fi
    done
    
    log_info "数据库迁移完成"
}

# =============================================================================
# Seed Functions
# =============================================================================

run_seeds() {
    log_info "执行种子数据导入..."
    
    local seeds_dir="${PROJECT_ROOT}/database/seeds"
    
    if [[ ! -d "${seeds_dir}" ]]; then
        log_warn "种子目录不存在: ${seeds_dir}"
        return 0
    fi
    
    # Get list of seed files sorted
    local seeds
    seeds=$(ls -1 "${seeds_dir}"/*.sql 2>/dev/null | sort)
    
    if [[ -z "$seeds" ]]; then
        log_warn "没有找到种子文件"
        return 0
    fi
    
    for seed in $seeds; do
        local filename
        filename=$(basename "$seed")
        log_info "执行种子: ${filename}"
        
        if [[ -n "${POSTGRES_PASSWORD}" ]]; then
            PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${POSTGRES_USER}" -d "${DB_NAME}" -f "$seed" 2>/dev/null || {
                log_warn "种子执行警告: ${filename} (可能已存在数据)"
            }
        else
            sudo -u postgres psql -d "${DB_NAME}" -f "$seed" || {
                log_warn "种子执行警告: ${filename} (可能已存在数据)"
            }
        fi
    done
    
    log_info "种子数据导入完成"
}

# =============================================================================
# Verify Functions
# =============================================================================

verify_setup() {
    log_info "验证数据库设置..."
    
    local test_result
    if [[ -n "${POSTGRES_PASSWORD}" ]]; then
        test_result=$(PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -tAc "SELECT current_database();" 2>/dev/null || echo "")
    else
        test_result=$(sudo -u postgres psql -d "${DB_NAME}" -tAc "SELECT current_database();" 2>/dev/null || echo "")
    fi
    
    if [[ "$test_result" == "${DB_NAME}" ]]; then
        log_info "✓ 数据库连接验证成功"
    else
        log_error "✗ 数据库连接验证失败"
        return 1
    fi
    
    # Check extensions
    local extensions
    if [[ -n "${POSTGRES_PASSWORD}" ]]; then
        extensions=$(PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${POSTGRES_USER}" -d "${DB_NAME}" -tAc "SELECT extname FROM pg_extension WHERE extname IN ('uuid-ossp', 'pgcrypto', 'pg_trgm');" 2>/dev/null)
    else
        extensions=$(sudo -u postgres psql -d "${DB_NAME}" -tAc "SELECT extname FROM pg_extension WHERE extname IN ('uuid-ossp', 'pgcrypto', 'pg_trgm');")
    fi
    
    if echo "$extensions" | grep -q "uuid-ossp"; then
        log_info "✓ uuid-ossp 扩展已启用"
    else
        log_warn "✗ uuid-ossp 扩展未找到"
    fi
    
    log_info "数据库初始化完成！"
}

# =============================================================================
# Main Function
# =============================================================================

show_help() {
    cat << EOF
RDP Database Initialization Script

用法: $0 [选项]

选项:
    -h, --help          显示帮助信息
    -i, --install       安装 PostgreSQL（如果不存在）
    -d, --database      指定数据库名称（默认: rdp）
    -u, --user          指定数据库用户（默认: rdp_user）
    -p, --password      指定数据库密码
    -H, --host          指定数据库主机（默认: localhost）
    -P, --port          指定数据库端口（默认: 5432）
    -s, --skip-seed     跳过种子数据导入
    -v, --verbose       显示详细输出

环境变量:
    RDP_DB_PASSWORD     数据库密码
    RDP_DB_HOST         数据库主机
    RDP_DB_PORT         数据库端口
    POSTGRES_USER       PostgreSQL管理员用户（默认: postgres）
    POSTGRES_PASSWORD   PostgreSQL管理员密码
    RDP_LOG_DIR         日志目录（默认: /var/log/rdp）

示例:
    # 基本初始化
    sudo $0

    # 指定密码
    sudo RDP_DB_PASSWORD=mypassword $0

    # 远程数据库
    $0 -H db.example.com -P 5432 -p mypassword
EOF
}

main() {
    local skip_seed=false
    local install_postgres=false
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -i|--install)
                install_postgres=true
                shift
                ;;
            -d|--database)
                DB_NAME="$2"
                shift 2
                ;;
            -u|--user)
                DB_USER="$2"
                shift 2
                ;;
            -p|--password)
                DB_PASSWORD="$2"
                shift 2
                ;;
            -H|--host)
                DB_HOST="$2"
                shift 2
                ;;
            -P|--port)
                DB_PORT="$2"
                shift 2
                ;;
            -s|--skip-seed)
                skip_seed=true
                shift
                ;;
            -v|--verbose)
                set -x
                shift
                ;;
            *)
                log_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # Print banner
    echo -e "${BLUE}"
    echo "╔════════════════════════════════════════════════════════════╗"
    echo "║     RDP Database Initialization Script                     ║"
    echo "║     Version: 1.0                                           ║"
    echo "╚════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    log_info "日志文件: ${LOG_FILE}"
    log_info "配置信息:"
    log_info "  数据库: ${DB_NAME}"
    log_info "  用户: ${DB_USER}"
    log_info "  主机: ${DB_HOST}:${DB_PORT}"
    
    # Check if running as root for local installation
    if [[ "${DB_HOST}" == "localhost" && "$EUID" -ne 0 && -z "${POSTGRES_PASSWORD}" ]]; then
        log_warn "建议以 root 用户运行此脚本，或使用 POSTGRES_PASSWORD 环境变量"
    fi
    
    # Step 1: Check/Install PostgreSQL
    if [[ "$install_postgres" == true ]]; then
        if ! check_postgres_installed; then
            install_postgresql
        fi
    else
        if ! check_postgres_installed; then
            log_error "PostgreSQL 未安装。使用 -i 选项自动安装，或手动安装后重试。"
            exit 1
        fi
    fi
    
    # Step 2: Check if PostgreSQL is running
    if ! check_postgres_running; then
        log_info "尝试启动 PostgreSQL 服务..."
        local service_name
        service_name=$(detect_postgres_service)
        systemctl start "${service_name}" || true
        
        sleep 2
        
        if ! check_postgres_running; then
            log_error "无法启动 PostgreSQL 服务"
            exit 1
        fi
    fi
    
    # Step 3: Create database
    create_database
    
    # Step 4: Create user
    create_user
    
    # Step 5: Create extensions
    create_extensions
    
    # Step 6: Run migrations
    run_migrations
    
    # Step 7: Run seeds (if not skipped)
    if [[ "$skip_seed" == false ]]; then
        run_seeds
    else
        log_info "跳过种子数据导入"
    fi
    
    # Step 8: Verify setup
    verify_setup
    
    # Print summary
    echo -e "${GREEN}"
    echo "╔════════════════════════════════════════════════════════════╗"
    echo "║     数据库初始化完成！                                      ║"
    echo "╠════════════════════════════════════════════════════════════╣"
    echo "║  数据库: ${DB_NAME}"
    echo "║  用户: ${DB_USER}"
    echo "║  密码: ${DB_PASSWORD}"
    echo "║  主机: ${DB_HOST}:${DB_PORT}"
    echo "╠════════════════════════════════════════════════════════════╣"
    echo "║  连接字符串:                                                ║"
    echo "║  postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}"
    echo "╚════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    log_info "详细日志保存在: ${LOG_FILE}"
}

# Run main function
main "$@"
