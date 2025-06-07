# Variáveis
APP_NAME=bin/app
MAIN=cmd/api/*.go
ENV=.env
API_DIR=cmd/api
SWAGGER_OUT=docs

# Targets
run: build
	@echo "Starting the Go server..."
	$(APP_NAME) -config $(ENV)

build:
	@mkdir -p bin
	@echo "Building the Go application..."
	GOOS=linux GOARCH=amd64 go build -o $(APP_NAME) $(MAIN)

up:
	@echo "Up Postgres with Docker Compose..."
	docker-compose up -d

down:
	@echo "Stopping Docker containers..."
	docker-compose down

deps:
	@echo "Downloading dependencies..."
	go mod tidy

clean:
	@echo "Clen up..."
	go clean
	rm -f $(APP_NAME)

.PHONY: swag
swag:
	cd $(API_DIR) && swag init --output ../../$(SWAGGER_OUT)

help:
	@echo ""
	@echo "Comandos disponíveis:"
	@echo "  make run       → roda o servidor Go"
	@echo "  make up        → sobe o postgres (docker-compose)"
	@echo "  make down      → derruba os containers"
	@echo "  make deps      → organiza módulos"
	@echo "  make clean     → limpa binários"
	@echo "  make swag      → gera documentação Swagger"
