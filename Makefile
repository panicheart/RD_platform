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
