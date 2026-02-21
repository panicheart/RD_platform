.PHONY: all install dev build lint test clean

all: install

install:
	cd apps/web && npm install
	cd services/api && go mod download

dev-frontend:
	cd apps/web && npm run dev

dev-backend:
	cd services/api && go run main.go

dev:
	@echo "Run 'make dev-frontend' and 'make dev-backend' in separate terminals"

build-frontend:
	cd apps/web && npm run build

build-backend:
	cd services/api && go build -o bin/api main.go

build: build-frontend build-backend

lint-frontend:
	cd apps/web && npm run lint

lint-backend:
	cd services/api && golangci-lint run ./...

lint: lint-frontend lint-backend

format-frontend:
	cd apps/web && npm run format

format-backend:
	cd services/api && gofmt -w . && goimports -w .

format: format-frontend format-backend

test-frontend:
	cd apps/web && npm run test

test-frontend-watch:
	cd apps/web && npm run test -- --watch

test-frontend-coverage:
	cd apps/web && npm run test:coverage

test-frontend-ui:
	cd apps/web && npm run test:ui

test-backend:
	cd services/api && go test -v ./...

test-backend-coverage:
	cd services/api && go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

test-backend-coverage-html:
	cd services/api && go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

test: test-backend test-frontend

test-coverage: test-backend-coverage test-frontend-coverage

clean:
	rm -rf apps/web/dist apps/web/node_modules
	rm -rf services/api/bin
	rm -rf tmp/

docker-build:
	docker build -t rdp-platform:latest .

deploy:
	cd deploy/scripts && sudo ./install.sh

db-migrate:
	@echo "Run database migrations"
	@echo "TODO: implement migration command"

db-seed:
	@echo "Run database seeds"
	@echo "TODO: implement seed command"

# =============================================================================
# 5-Agent Team 管理命令
# =============================================================================

.PHONY: agent-team agent-team-start agent-team-status agent-team-init

# 5-Agent Team 帮助
agent-team:
	@echo "RDP 5-Agent Team 管理命令"
	@echo ""
	@echo "可用命令:"
	@echo "  make agent-team-start    - 启动5个Agent会话"
	@echo "  make agent-team-status   - 查看任务看板"
	@echo "  make agent-team-init     - 初始化Phase 1任务"
	@echo "  make agent-pm            - 启动 PM-Agent"
	@echo "  make agent-architect     - 启动 Architect-Agent"
	@echo "  make agent-backend       - 启动 Backend-Agent"
	@echo "  make agent-frontend      - 启动 Frontend-Agent"
	@echo "  make agent-devops        - 启动 DevOps-Agent"

# 初始化5-Agent任务
agent-team-init:
	@echo "🚀 初始化5-Agent Team任务..."
	python3 agents/5-agent-team/coordinator.py init

# 查看任务看板
agent-team-status:
	@echo "📋 5-Agent Team 任务看板"
	python3 agents/5-agent-team/coordinator.py status

# 启动所有Agent会话说明
agent-team-start:
	@echo "🚀 启动5-Agent Team"
	@echo ""
	@echo "请在5个不同终端分别执行以下命令:"
	@echo ""
	@echo "终端1 (PM-Agent):"
	@echo "  make agent-pm"
	@echo ""
	@echo "终端2 (Architect-Agent):"
	@echo "  make agent-architect"
	@echo ""
	@echo "终端3 (Backend-Agent):"
	@echo "  make agent-backend"
	@echo ""
	@echo "终端4 (Frontend-Agent):"
	@echo "  make agent-frontend"
	@echo ""
	@echo "终端5 (DevOps-Agent):"
	@echo "  make agent-devops"
	@echo ""
	@echo "或者使用脚本:"
	@echo "  ./agents/5-agent-team/start-pm.sh"
	@echo "  ./agents/5-agent-team/start-architect.sh"
	@echo "  ./agents/5-agent-team/start-backend.sh"
	@echo "  ./agents/5-agent-team/start-frontend.sh"
	@echo "  ./agents/5-agent-team/start-devops.sh"

# 一键启动所有Agent
agent-team-start-all:
	@echo "🚀 一键启动5-Agent Team..."
	./agents/5-agent-team/start-all.sh

# 停止所有Agent
agent-team-stop:
	@echo "🛑 停止5-Agent Team..."
	./agents/5-agent-team/stop-all.sh

# 启动单个Agent
agent-pm:
	@echo "🚀 启动 PM-Agent (项目经理)..."
	./agents/5-agent-team/start-pm.sh

agent-architect:
	@echo "🚀 启动 Architect-Agent (架构师)..."
	./agents/5-agent-team/start-architect.sh

agent-backend:
	@echo "🚀 启动 Backend-Agent (后端开发)..."
	./agents/5-agent-team/start-backend.sh

agent-frontend:
	@echo "🚀 启动 Frontend-Agent (前端开发)..."
	./agents/5-agent-team/start-frontend.sh

agent-devops:
	@echo "🚀 启动 DevOps-Agent (运维部署)..."
	./agents/5-agent-team/start-devops.sh
