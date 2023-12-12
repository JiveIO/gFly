.PHONY: clean critic security lint test build run

APP_NAME = app
CLI_NAME = artisan
BUILD_DIR = $(PWD)/build
MIGRATION_FOLDER = $(PWD)/database/migrations/postgresql
DATABASE_URL = postgres://vinh:@localhost:5432/gfly?sslmode=disable

all: clean critic security lint test swag build

clean:
	rm -rf ./build

critic:
	gocritic check -enableAll -disable=unnamedResult,unlabelStmt,hugeParam,singleCaseSwitch ./...

security:
	gosec ./...

lint:
	golangci-lint run ./...

test:
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

test.cover:
	go tool cover -html=cover.out

build: clean critic security lint test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(CLI_NAME) app/console/cli.go

run: swag build
	$(BUILD_DIR)/$(APP_NAME)

start: run

migrate.up:
	migrate -path $(MIGRATION_FOLDER) -database "$(DATABASE_URL)" up

migrate.down:
	migrate -path $(MIGRATION_FOLDER) -database "$(DATABASE_URL)" down

air:
	air main.go

swag:
	swag init
	cp ./docs/swagger.json ./public/docs

release:
	mkdir -p bin
	# 64 bit - Windows
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-amd64.exe *.go
	# 64-bit - Mac
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-amd64-darwin *.go
	# 64-bit - Mac ARM
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-arm64-darwin *.go
	# 64-bit - Linux
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-amd64-linux *.go
	# 64-bit - Linux ARM
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-arm64-linux *.go

docker.image:
	docker-compose -f docker/docker-compose.yml build --no-cache web

docker.start: docker.run docker.logs
docker.checking: docker.critic docker.security docker.lint

docker.migrate.up:
	docker exec -it gfly-web migrate -path /app/database/migrations/mysql -database "mysql://user:secret@tcp(db:3306)/gfly" up

docker.migrate.down:
	docker exec -it gfly-web migrate -path /app/database/migrations/mysql -database "mysql://user:secret@tcp(db:3306)/gfly" down

docker.critic:
	docker exec -it gfly-web gocritic check -enableAll -disable=unnamedResult,unlabelStmt,hugeParam,singleCaseSwitch ./...

docker.security:
	docker exec -it gfly-web gosec ./...

docker.lint:
	docker exec -it gfly-web golangci-lint run ./...

docker.test:
	docker exec -it gfly-web go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	docker exec -it gfly-web go tool cover -func=cover.out

docker.build: docker.checking docker.test
	docker exec -it gfly-web build_app
	docker exec -it gfly-web build_artisan

docker.release:
	docker exec -it gfly-web release_windows_64
	docker exec -it gfly-web release_mac_amd64
	docker exec -it gfly-web release_mac_arm64
	docker exec -it gfly-web release_linux_amd64
	docker exec -it gfly-web release_linux_arm64

docker.swag:
	docker exec -it gfly-web swag init
	docker exec -it gfly-web cp /app/docs/swagger.json /app/public/docs

docker.run:
	docker-compose -f docker/docker-compose.yml -p gfly up -d web

docker.logs:
	docker-compose -f docker/docker-compose.yml -p gfly logs -f web

docker.shell:
	docker-compose -f docker/docker-compose.yml -p gfly exec web bash

docker.stop:
	docker-compose -f docker/docker-compose.yml -p gfly kill

docker.destroy:
	docker-compose -f docker/docker-compose.yml -p gfly down

docker.drop: docker.stop docker.destroy