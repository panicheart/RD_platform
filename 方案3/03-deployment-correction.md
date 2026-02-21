# 部署方案修正文档

> **文档编号**: RDP-DEPLOY-CORRECTION-001  
> **版本**: V1.0  
> **编制日期**: 2026年2月21日  
> **密级**: 内部公开

---

## 1. 变更说明

### 1.1 变更原因

根据原始需求文档的明确要求：

> "**项目不使用 docker 和 k8s 这类虚拟技术**"

原《02-详细实施方案.pdf》中采用的 **Docker Compose** 部署方式与该约束冲突，现修正为**裸机部署（systemd服务）**方案。

### 1.2 变更内容

| 项目 | 原方案 | 修正后方案 |
|------|--------|------------|
| **部署方式** | Docker Compose 容器编排 | 裸机二进制部署 + systemd服务管理 |
| **服务隔离** | 容器隔离 | 进程级隔离 + 独立用户权限 |
| **资源管理** | Docker资源限制 | systemd资源限制 (MemoryMax/CPUQuota) |
| **日志管理** | Docker日志驱动 | journald + rsyslog |
| **服务发现** | Docker网络/DNS | 本地端口绑定 + Nginx反向代理 |

### 1.3 受影响的文档内容

| 文档 | 页码/章节 | 原内容 | 修正后 |
|------|----------|--------|--------|
| 03-需求规格说明书 | 第2页，1.1节 | "部署环境：离线局域网，**Docker Compose单节点**" | "部署环境：离线局域网，**裸机部署(systemd服务)**" |
| 03-需求规格说明书 | 第4页，2.1节 | "容器化：**Docker + Docker Compose** 24.x [MUST]" | "服务管理：**systemd服务 + 裸机二进制部署** [MUST]" |
| 02-详细实施方案 | 第4页，2.1节 | "基础设施层 — **Docker Compose**" | "基础设施层 — **systemd服务管理**" |
| 02-详细实施方案 | 第7页，3.1节 | "**Docker + Docker Compose** 24.x" | "**systemd + 裸机部署**" |
| 02-详细实施方案 | 第20页，6.1节 | "**Docker Compose** 服务集群" | "**systemd服务**集群" |
| 02-详细实施方案 | 第22页，7.1节 | "**Docker部署方案**" | "**裸机部署方案**" |
| README.md | 技术栈 | "部署：**Docker Compose**" | "部署：**裸机部署（systemd服务）**" |

---

## 2. 裸机部署架构

### 2.1 系统架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                     设计师工作站 (浏览器 + 辅助程序)              │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Nginx 反向代理 (1.25+)                       │
│           TLS终止 | 路由分发 | 静态资源 | 负载均衡                 │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
          ┌───────────────────────┼───────────────────────┐
          │                       │                       │
          ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   rdp-portal    │    │   rdp-api       │    │   rdp-dev       │
│   (前端Shell)   │    │   (Go API)      │    │   (项目开发)    │
│   :3000         │    │   :8080         │    │   :3001         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
          │                       │                       │
          └───────────────────────┼───────────────────────┘
                                  │
          ┌───────────────────────┼───────────────────────┐
          │                       │                       │
          ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Casdoor       │    │   Gitea         │    │   Mattermost    │
│   (认证/RBAC)   │    │   (Git服务)     │    │   (即时通讯)    │
│   :8000         │    │   :3002         │    │   :8065         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
          │                       │                       │
          └───────────────────────┼───────────────────────┘
                                  │
          ┌───────────────────────┼───────────────────────┐
          │                       │                       │
          ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │   Redis         │    │   MinIO         │
│   (主数据库)    │    │   (缓存/会话)   │    │   (对象存储)    │
│   :5432         │    │   :6379         │    │   :9000/9001    │
└─────────────────┘    └─────────────────┘    └─────────────────┘

辅助服务:
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   MeiliSearch   │    │   Wiki.js       │    │   备份服务      │
│   (搜索引擎)    │    │   (知识库)      │    │   (定时任务)    │
│   :7700         │    │   :3003         │    │   systemd timer │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 2.2 部署目录结构

```
/opt/rdp/                          # 主部署目录
├── bin/                           # 可执行文件
│   ├── rdp-api                    # 主API服务
│   ├── rdp-portal/                # 前端静态文件
│   ├── rdp-dev/                   # 项目开发模块前端
│   └── rdp-shelf/                 # 货架模块前端
├── config/                        # 配置文件
│   ├── nginx/
│   │   ├── nginx.conf
│   │   └── sites-available/
│   ├── systemd/
│   │   ├── rdp-api.service
│   │   ├── rdp-portal.service
│   │   └── ...
│   ├── casdoor/
│   ├── gitea/
│   └── mattermost/
├── data/                          # 数据目录
│   ├── postgresql/                # PostgreSQL数据
│   ├── redis/                     # Redis数据
│   ├── minio/                     # MinIO存储
│   ├── gitea/                     # Gitea仓库
│   ├── mattermost/                # Mattermost数据
│   ├── meilisearch/               # 搜索索引
│   └── backups/                   # 备份文件
├── logs/                          # 日志目录
│   ├── rdp-api/
│   ├── nginx/
│   └── ...
├── scripts/                       # 运维脚本
│   ├── install.sh                 # 安装脚本
│   ├── backup.sh                  # 备份脚本
│   ├── update.sh                  # 更新脚本
│   └── health-check.sh            # 健康检查
└── temp/                          # 临时文件
```

---

## 3. 服务部署清单

### 3.1 基础依赖安装

```bash
# Ubuntu Server 22.04 LTS / CentOS 8+ / 麒麟OS

# 1. 系统更新
sudo apt update && sudo apt upgrade -y

# 2. 安装基础工具
sudo apt install -y wget curl vim htop net-tools \
    build-essential git unzip tar

# 3. 创建独立用户
sudo useradd -r -s /bin/false rdp-user
sudo mkdir -p /opt/rdp
sudo chown -R rdp-user:rdp-user /opt/rdp
```

### 3.2 数据库服务

#### PostgreSQL 16

```bash
# 安装 PostgreSQL 16
sudo apt install -y postgresql-16 postgresql-client-16

# 配置数据目录到 /opt/rdp/data/postgresql
sudo systemctl stop postgresql
sudo mkdir -p /opt/rdp/data/postgresql
sudo chown -R postgres:postgres /opt/rdp/data/postgresql
sudo -u postgres initdb -D /opt/rdp/data/postgresql

# 修改 systemd 服务使用新数据目录
sudo systemctl edit postgresql
# 添加：
# [Service]
# Environment=PGDATA=/opt/rdp/data/postgresql

sudo systemctl daemon-reload
sudo systemctl enable --now postgresql

# 创建数据库和用户
sudo -u postgres psql -c "CREATE USER rdp WITH PASSWORD 'your_password';"
sudo -u postgres psql -c "CREATE DATABASE rdp_db OWNER rdp;"
sudo -u postgres psql -c "CREATE DATABASE casdoor OWNER rdp;"
sudo -u postgres psql -c "CREATE DATABASE gitea OWNER rdp;"
sudo -u postgres psql -c "CREATE DATABASE mattermost OWNER rdp;"
```

#### Redis 7

```bash
# 安装 Redis
sudo apt install -y redis-server

# 配置
sudo tee /etc/redis/redis.conf > /dev/null <<'EOF'
bind 127.0.0.1
port 6379
dir /opt/rdp/data/redis
maxmemory 2gb
maxmemory-policy allkeys-lru
requirepass your_redis_password
EOF

sudo mkdir -p /opt/rdp/data/redis
sudo chown -R redis:redis /opt/rdp/data/redis
sudo systemctl restart redis-server
sudo systemctl enable redis-server
```

### 3.3 应用服务

#### 主 API 服务 (Go)

```bash
# 创建 systemd 服务
sudo tee /etc/systemd/system/rdp-api.service > /dev/null <<'EOF'
[Unit]
Description=RDP Platform API Service
After=network.target postgresql.service redis.service
Wants=postgresql.service redis.service

[Service]
Type=simple
User=rdp-user
Group=rdp-user
WorkingDirectory=/opt/rdp
ExecStart=/opt/rdp/bin/rdp-api
Restart=always
RestartSec=5

# 资源限制
MemoryMax=2G
CPUQuota=200%

# 环境变量
Environment="DB_HOST=127.0.0.1"
Environment="DB_PORT=5432"
Environment="DB_USER=rdp"
Environment="DB_PASSWORD=your_password"
Environment="DB_NAME=rdp_db"
Environment="REDIS_HOST=127.0.0.1"
Environment="REDIS_PORT=6379"
Environment="REDIS_PASSWORD=your_redis_password"
Environment="PORT=8080"

# 日志
StandardOutput=journal
StandardError=journal
SyslogIdentifier=rdp-api

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable rdp-api.service
```

#### 前端服务 (Nginx 托管)

```bash
# 安装 Nginx
sudo apt install -y nginx

# 配置前端目录
sudo mkdir -p /opt/rdp/bin/rdp-portal
sudo mkdir -p /opt/rdp/bin/rdp-dev
sudo mkdir -p /opt/rdp/bin/rdp-shelf

# 创建 systemd 服务（用于开发服务器模式）
sudo tee /etc/systemd/system/rdp-portal.service > /dev/null <<'EOF'
[Unit]
Description=RDP Portal Frontend
After=network.target

[Service]
Type=simple
User=rdp-user
Group=rdp-user
WorkingDirectory=/opt/rdp/bin/rdp-portal
ExecStart=/usr/bin/python3 -m http.server 3000
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 生产环境使用 Nginx 直接托管静态文件
sudo tee /etc/nginx/sites-available/rdp-portal > /dev/null <<'EOF'
server {
    listen 3000;
    server_name localhost;
    
    root /opt/rdp/bin/rdp-portal/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
EOF

sudo ln -s /etc/nginx/sites-available/rdp-portal /etc/nginx/sites-enabled/
sudo systemctl restart nginx
```

### 3.4 开源组件服务

#### Casdoor (二进制部署)

```bash
# 下载 Casdoor 二进制
wget https://github.com/casdoor/casdoor/releases/download/v1.700.0/casdoor-v1.700.0-linux-amd64.tar.gz
tar -xzf casdoor-v1.700.0-linux-amd64.tar.gz
sudo mv casdoor /opt/rdp/bin/
sudo chmod +x /opt/rdp/bin/casdoor

# 配置文件
sudo mkdir -p /opt/rdp/config/casdoor
sudo tee /opt/rdp/config/casdoor/app.conf > /dev/null <<'EOF'
appname = casdoor
httpport = 8000
runmode = prod

[database]
adapter = postgres
host = 127.0.0.1
port = 5432
user = rdp
password = your_password
database = casdoor
EOF

# systemd 服务
sudo tee /etc/systemd/system/casdoor.service > /dev/null <<'EOF'
[Unit]
Description=Casdoor Identity Platform
After=network.target postgresql.service

[Service]
Type=simple
User=rdp-user
Group=rdp-user
WorkingDirectory=/opt/rdp/config/casdoor
ExecStart=/opt/rdp/bin/casdoor
Restart=always
RestartSec=5
MemoryMax=1G

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable casdoor.service
```

#### Gitea (二进制部署)

```bash
# 下载 Gitea
wget -O gitea https://dl.gitea.com/gitea/1.22.0/gitea-1.22.0-linux-amd64
sudo mv gitea /opt/rdp/bin/
sudo chmod +x /opt/rdp/bin/gitea

# 初始化配置
sudo mkdir -p /opt/rdp/data/gitea
sudo chown -R rdp-user:rdp-user /opt/rdp/data/gitea

# systemd 服务
sudo tee /etc/systemd/system/gitea.service > /dev/null <<'EOF'
[Unit]
Description=Gitea Git Service
After=network.target postgresql.service

[Service]
Type=simple
User=rdp-user
Group=rdp-user
WorkingDirectory=/opt/rdp/data/gitea
ExecStart=/opt/rdp/bin/gitea web -c /opt/rdp/config/gitea/app.ini
Restart=always
RestartSec=5
MemoryMax=2G

[Install]
WantedBy=multi-user.target
EOF
```

#### MinIO (二进制部署)

```bash
# 下载 MinIO
wget https://dl.min.io/server/minio/release/linux-amd64/minio
sudo mv minio /opt/rdp/bin/
sudo chmod +x /opt/rdp/bin/minio

# 数据目录
sudo mkdir -p /opt/rdp/data/minio
sudo chown -R rdp-user:rdp-user /opt/rdp/data/minio

# systemd 服务
sudo tee /etc/systemd/system/minio.service > /dev/null <<'EOF'
[Unit]
Description=MinIO Object Storage
After=network.target

[Service]
Type=simple
User=rdp-user
Group=rdp-user
WorkingDirectory=/opt/rdp/data/minio
Environment="MINIO_ROOT_USER=rdp-admin"
Environment="MINIO_ROOT_PASSWORD=your_minio_password"
ExecStart=/opt/rdp/bin/minio server /opt/rdp/data/minio --address :9000 --console-address :9001
Restart=always
RestartSec=5
MemoryMax=4G

[Install]
WantedBy=multi-user.target
EOF
```

#### Mattermost (二进制部署)

```bash
# 下载 Mattermost
wget https://releases.mattermost.com/9.11.0/mattermost-9.11.0-linux-amd64.tar.gz
tar -xzf mattermost-9.11.0-linux-amd64.tar.gz
sudo mv mattermost /opt/rdp/
sudo mkdir -p /opt/rdp/mattermost/data
sudo chown -R rdp-user:rdp-user /opt/rdp/mattermost

# systemd 服务
sudo tee /etc/systemd/system/mattermost.service > /dev/null <<'EOF'
[Unit]
Description=Mattermost Team Collaboration
After=network.target postgresql.service

[Service]
Type=notify
User=rdp-user
Group=rdp-user
WorkingDirectory=/opt/rdp/mattermost
ExecStart=/opt/rdp/mattermost/bin/mattermost
Restart=always
RestartSec=5
MemoryMax=2G

[Install]
WantedBy=multi-user.target
EOF
```

#### MeiliSearch (二进制部署)

```bash
# 下载 MeiliSearch
curl -L https://install.meilisearch.com | sh
sudo mv meilisearch /opt/rdp/bin/
sudo chmod +x /opt/rdp/bin/meilisearch

# 数据目录
sudo mkdir -p /opt/rdp/data/meilisearch
sudo chown -R rdp-user:rdp-user /opt/rdp/data/meilisearch

# systemd 服务
sudo tee /etc/systemd/system/meilisearch.service > /dev/null <<'EOF'
[Unit]
Description=MeiliSearch Engine
After=network.target

[Service]
Type=simple
User=rdp-user
Group=rdp-user
WorkingDirectory=/opt/rdp/data/meilisearch
ExecStart=/opt/rdp/bin/meilisearch --db-path /opt/rdp/data/meilisearch --http-addr 127.0.0.1:7700 --master-key your_meili_key
Restart=always
RestartSec=5
MemoryMax=1G

[Install]
WantedBy=multi-user.target
EOF
```

---

## 4. Nginx 统一配置

```bash
sudo tee /etc/nginx/nginx.conf > /dev/null <<'EOF'
user www-data;
worker_processes auto;
pid /run/nginx.pid;

events {
    worker_connections 4096;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # 日志格式
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /opt/rdp/logs/nginx/access.log main;
    error_log /opt/rdp/logs/nginx/error.log warn;

    # 性能优化
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    client_max_body_size 100M;

    # Gzip压缩
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml;

    # 上游服务定义
    upstream rdp_api {
        server 127.0.0.1:8080;
    }

    upstream casdoor {
        server 127.0.0.1:8000;
    }

    upstream gitea {
        server 127.0.0.1:3002;
    }

    upstream mattermost {
        server 127.0.0.1:8065;
    }

    upstream minio {
        server 127.0.0.1:9000;
    }

    upstream minio_console {
        server 127.0.0.1:9001;
    }

    # 主站点配置
    server {
        listen 80;
        server_name _;
        
        # 重定向到 HTTPS (如果配置了SSL)
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name rdp.local;

        # SSL证书配置
        ssl_certificate /opt/rdp/config/nginx/ssl/rdp.crt;
        ssl_certificate_key /opt/rdp/config/nginx/ssl/rdp.key;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers HIGH:!aNULL:!MD5;

        # 前端静态资源
        location / {
            root /opt/rdp/bin/rdp-portal/dist;
            try_files $uri $uri/ /index.html;
            expires 1d;
        }

        # API代理
        location /api/ {
            proxy_pass http://rdp_api/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 60s;
        }

        # Casdoor认证
        location /auth/ {
            proxy_pass http://casdoor/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }

        # Gitea Git服务
        location /git/ {
            proxy_pass http://gitea/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }

        # Mattermost即时通讯
        location /chat/ {
            proxy_pass http://mattermost/;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_set_header Host $host;
        }

        # MinIO对象存储
        location /storage/ {
            proxy_pass http://minio/;
            proxy_set_header Host $host;
        }

        # 文件上传
        location /upload/ {
            proxy_pass http://minio/;
            client_max_body_size 1G;
        }
    }
}
EOF

sudo mkdir -p /opt/rdp/logs/nginx /opt/rdp/config/nginx/ssl
sudo nginx -t && sudo systemctl restart nginx
```

---

## 5. 运维脚本

### 5.1 一键安装脚本

```bash
#!/bin/bash
# /opt/rdp/scripts/install.sh

set -e

echo "=========================================="
echo "RDP Platform - Bare Metal Installation"
echo "=========================================="

# 1. 系统检查
if [[ $EUID -ne 0 ]]; then
   echo "请使用 sudo 运行此脚本"
   exit 1
fi

# 2. 创建目录结构
echo "[1/5] 创建目录结构..."
mkdir -p /opt/rdp/{bin,config,data,logs,scripts,temp}
mkdir -p /opt/rdp/data/{postgresql,redis,minio,gitea,mattermost,meilisearch,backups}
mkdir -p /opt/rdp/logs/{nginx,rdp-api}
mkdir -p /opt/rdp/config/{nginx,casdoor,gitea,mattermost}

# 3. 创建用户
echo "[2/5] 创建服务用户..."
id -u rdp-user &>/dev/null || useradd -r -s /bin/false rdp-user
chown -R rdp-user:rdp-user /opt/rdp

# 4. 安装基础依赖
echo "[3/5] 安装基础依赖..."
apt update
apt install -y wget curl vim htop net-tools build-essential git unzip tar nginx

# 5. 安装数据库
echo "[4/5] 安装 PostgreSQL 和 Redis..."
apt install -y postgresql-16 postgresql-client-16 redis-server

# 6. 启动服务
echo "[5/5] 配置 systemd 服务..."
systemctl daemon-reload
systemctl enable postgresql redis-server nginx

echo "=========================================="
echo "基础环境安装完成"
echo "下一步: 配置数据库并安装应用服务"
echo "=========================================="
```

### 5.2 备份脚本

```bash
#!/bin/bash
# /opt/rdp/scripts/backup.sh

BACKUP_DIR="/opt/rdp/data/backups/$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

# 备份 PostgreSQL
echo "备份 PostgreSQL..."
sudo -u postgres pg_dump rdp_db > "$BACKUP_DIR/rdp_db.sql"
sudo -u postgres pg_dump casdoor > "$BACKUP_DIR/casdoor.sql"
sudo -u postgres pg_dump gitea > "$BACKUP_DIR/gitea.sql"

# 备份 Redis
echo "备份 Redis..."
redis-cli BGSAVE
cp /opt/rdp/data/redis/dump.rdb "$BACKUP_DIR/redis.rdb"

# 备份 MinIO
echo "备份 MinIO..."
mc mirror /opt/rdp/data/minio "$BACKUP_DIR/minio"

# 备份 Gitea
echo "备份 Gitea..."
tar -czf "$BACKUP_DIR/gitea.tar.gz" -C /opt/rdp/data/gitea .

# 备份配置文件
echo "备份配置文件..."
tar -czf "$BACKUP_DIR/config.tar.gz" -C /opt/rdp/config .

# 清理旧备份 (保留7天)
find /opt/rdp/data/backups -type d -mtime +7 -exec rm -rf {} + 2>/dev/null

echo "备份完成: $BACKUP_DIR"
```

### 5.3 健康检查脚本

```bash
#!/bin/bash
# /opt/rdp/scripts/health-check.sh

check_service() {
    local name=$1
    local port=$2
    if nc -z 127.0.0.1 $port 2>/dev/null; then
        echo "✅ $name (port $port) - 正常"
        return 0
    else
        echo "❌ $name (port $port) - 异常"
        return 1
    fi
}

echo "=========================================="
echo "RDP Platform Health Check"
echo "=========================================="
echo ""

check_service "PostgreSQL" 5432
check_service "Redis" 6379
check_service "RDP API" 8080
check_service "Casdoor" 8000
check_service "Gitea" 3002
check_service "Mattermost" 8065
check_service "MinIO" 9000
check_service "MeiliSearch" 7700
check_service "Nginx" 80
check_service "Nginx SSL" 443

echo ""
echo "=========================================="
echo "磁盘空间:"
df -h /opt/rdp
echo ""
echo "内存使用:"
free -h
echo ""
echo "=========================================="
```

### 5.4 服务管理命令速查

```bash
# 查看所有服务状态
systemctl status postgresql redis-server nginx rdp-api casdoor gitea mattermost minio meilisearch

# 启动所有服务
systemctl start postgresql redis-server nginx rdp-api casdoor gitea mattermost minio meilisearch

# 停止所有服务
systemctl stop rdp-api casdoor gitea mattermost minio meilisearch nginx

# 查看日志
journalctl -u rdp-api -f          # 实时查看API日志
journalctl -u casdoor --since "1 hour ago"  # 查看Casdoor最近1小时日志

# 重启单个服务
systemctl restart rdp-api
```

---

## 6. 与 Docker 方案的对比优势

| 优势项 | 说明 |
|--------|------|
| **符合原始需求** | 完全遵守"不使用Docker"的约束 |
| **资源利用率** | 无容器开销，内存/CPU使用更直接 |
| **性能** | 无容器化层，IO和网络性能更优 |
| **安全性** | 传统Linux权限模型，审计更直接 |
| **兼容性** | 与现有运维体系（监控、备份）无缝集成 |
| **可控性** | 每个服务独立管理，故障隔离清晰 |

---

## 7. 附录

### 7.1 端口分配表

| 服务 | 端口 | 说明 |
|------|------|------|
| Nginx HTTP | 80 | HTTP入口 |
| Nginx HTTPS | 443 | HTTPS入口 |
| RDP Portal | 3000 | 前端Shell（开发模式） |
| RDP Dev | 3001 | 项目开发模块 |
| Gitea | 3002 | Git服务 |
| Wiki.js | 3003 | 知识库（可选） |
| Casdoor | 8000 | 认证服务 |
| RDP API | 8080 | 主API服务 |
| Mattermost | 8065 | 即时通讯 |
| PostgreSQL | 5432 | 主数据库 |
| Redis | 6379 | 缓存/会话 |
| MinIO API | 9000 | 对象存储API |
| MinIO Console | 9001 | 对象存储控制台 |
| MeiliSearch | 7700 | 搜索引擎 |

### 7.2 文档变更对照表

| 原文档 | 页码 | 原内容 | 修正后内容 |
|--------|------|--------|------------|
| 03-需求规格说明书 | 第2页 | 部署环境：离线局域网，**Docker Compose单节点** | 部署环境：离线局域网，**裸机部署(systemd服务)** |
| 03-需求规格说明书 | 第4页 | **Docker + Docker Compose** 24.x [MUST] | **systemd + 裸机二进制部署** [MUST] |
| 02-详细实施方案 | 第4页 | 基础设施层 — **Docker Compose** | 基础设施层 — **systemd服务管理** |
| 02-详细实施方案 | 第7页 | **Docker + Docker Compose** 24.x | **systemd + 裸机部署** |
| 02-详细实施方案 | 第20页 | **Docker Compose** 服务集群 | **systemd服务**集群 |
| 02-详细实施方案 | 第22页 | **Docker部署方案** | **裸机部署方案** |
| README.md | 技术栈 | **Docker Compose**（离线局域网） | **裸机部署（systemd服务）** |

---

**文档结束**
