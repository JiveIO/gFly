.PHONY: lint test build run dev

APP_NAME = app
CLI_NAME = artisan
BUILD_DIR = $(PWD)/build
MIGRATION_FOLDER = database/migrations/postgresql
DATABASE_URL = postgres://user:secret@localhost:5432/gfly?sslmode=disable

all: lint test doc build

lint:
	golangci-lint run ./...

test:
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out
	go tool cover -html=cover.out

build: lint test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(CLI_NAME) app/console/cli.go

run: lint test doc build
	$(BUILD_DIR)/$(APP_NAME)

start: run

migrate.up:
	migrate -path $(MIGRATION_FOLDER) -database "$(DATABASE_URL)" up

migrate.down:
	migrate -path $(MIGRATION_FOLDER) -database "$(DATABASE_URL)" down

dev:
	air -build.exclude_dir=node_modules,public,resources,Dev,bin,build,dist,docker,storage,tmp,database,docs main.go

clean:
	go mod tidy
	go clean -cache
	go clean -testcache

doc:
	swag init
	cp ./docs/swagger.json ./public/docs/

docker.run:
	docker-compose -f docker/docker-compose.yml -p gfly up -d db
	docker-compose -f docker/docker-compose.yml -p gfly up -d mail
	docker-compose -f docker/docker-compose.yml -p gfly up -d redis

docker.logs:
	docker-compose -f docker/docker-compose.yml -p gfly logs -f db &
	docker-compose -f docker/docker-compose.yml -p gfly logs -f mail &
	docker-compose -f docker/docker-compose.yml -p gfly logs -f redis &

docker.stop:
	docker-compose -f docker/docker-compose.yml -p gfly kill

docker.delete:
	docker-compose -f docker/docker-compose.yml -p gfly down
