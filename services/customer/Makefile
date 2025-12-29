.PHONY: migrate-up migrate-down migrate-create db-create db-drop test

# database config
DB_HOST ?= localhost
DB_PORT ?= 3306
DB_USER ?= root
DB_NAME ?= subscription-billing-sys
DB_PASSWORD ?= 

# migration paths
MIGRATION_PATH ?= file://migrations

# conn strings
# DB_DSN = mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?charset=utf8mb4&parseTime=True&loc=Local
DB_DSN = mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?multiStatements=true
GORM_DSN = $(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?charset=uft8mb4&parseTime=True&loc=Local


# create db
db-create:
	@echo "Creating database $(DB_NAME)..."
	@/opt/lampp/bin/mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASSWORD) -e "create database if not exists \`$(DB_NAME)\` character set utf8mb4 collate utf8mb4_unicode_ci"
	@echo "Database created successfully"


# run migrations up
migrate-up:
	@echo "Running migrations up..."
	@migrate -path ./migrations -database "$(DB_DSN)" up
	@echo "Migrations completed successfully"


# rollback migrations
migrate-down:
	@echo "Rolling back migrations"
	@migrate -path migrations -database "$(DB_DSN)" down 1
	@echo "Rollback completed"

# create new migration files
migrate-create:
	@echo "Creating a new migration file..."
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $${name}

# reset database [drop, create, migrate]
db-reset:
	$(MAKE) db-drop
	$(MAKE) db-create
	$(MAKE) migrate-up

# test db conn
test-db:
	@echo "Testing database connection..."
	@/opt/lampp/bin/mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASSWORD) -e "select 1;" $(DB_NAME) && echo "Connection successful"

# show migration status
migrate-status:
	@echo "Migration status:"
	@migrate -path ./migrations -database "$(DB_DSN)" version

# run the app
run:
	@go run cmd/api/main.go


# build the app
build:
	@go build -o /bin/subscription-billing-sys/cmd/api/main.go

# install dependencies
deps:
	@go mod download
	@go mod tidy

# run tests
test: 
	@go test ./... -v 

# clean build build artifacts
clean:
	@rm -rf bin/
	@go clean

