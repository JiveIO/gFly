.PHONY: critic security lint test build run dev

APP_NAME = app
CLI_NAME = artisan
BUILD_DIR = $(PWD)/build
DOCKER_PATH = /home/gfly/app
MIGRATION_FOLDER = $(DOCKER_PATH)/database/migrations/postgresql
DATABASE_URL = postgres://user:secret@db:5432/gfly?sslmode=disable
#MIGRATION_FOLDER = $(DOCKER_PATH)/database/migrations/mysql
#DATABASE_URL = mysql://user:secret@tcp(db:3306)/gfly

all: critic security lint test doc build

critic:
	gocritic check -enableAll -disable=unnamedResult,unlabelStmt,hugeParam,singleCaseSwitch,builtinShadow,typeAssertChain ./...

security:
	gosec -exclude-dir=core -exclude=G103 ./...

vulncheck:
	govulncheck ./...

lint:
	golangci-lint run ./...

test:
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

test.cover:
	go tool cover -html=cover.out

build: critic security vulncheck lint test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(CLI_NAME) app/console/cli.go

run: doc build
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

release:
	# 64-bit - Linux
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-amd64-linux *.go
	# 64-bit - Linux ARM
	#GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-arm64-linux *.go

release.windows:
	mkdir -p bin
	# 64 bit - Windows
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-amd64.exe *.go

release.mac:
	mkdir -p bin
	# 64-bit - Mac
	#GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-amd64-darwin *.go
	# 64-bit - Mac ARM
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-arm64-darwin *.go

docker.image:
	docker-compose -f docker/docker-compose.yml build --no-cache --build-arg hostUID=1000 --build-arg hostGID=1000 web

docker.start: docker.run docker.logs
docker.checking: docker.critic docker.vulncheck docker.security docker.lint

docker.migrate.up:
	docker exec -it --user gfly gfly-web migrate -path $(MIGRATION_FOLDER) -database "$(DATABASE_URL)" up

docker.migrate.down:
	docker exec -it --user gfly gfly-web migrate -path $(MIGRATION_FOLDER) -database "$(DATABASE_URL)" down 1

docker.critic:
	docker exec -it --user gfly gfly-web gocritic check -enableAll -disable=unnamedResult,unlabelStmt,hugeParam,singleCaseSwitch,builtinShadow,typeAssertChain ./...

docker.vulncheck:
	docker exec -it --user gfly gfly-web govulncheck ./...

docker.security:
	docker exec -it --user gfly gfly-web gosec -exclude=G103 ./...

docker.lint:
	docker exec -it --user gfly gfly-web golangci-lint run ./...

docker.test:
	docker exec -it --user gfly gfly-web go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	docker exec -it --user gfly gfly-web go tool cover -func=cover.out

docker.doc:
	docker exec -it --user gfly gfly-web swag init
	docker exec -it --user gfly gfly-web cp $(DOCKER_PATH)/docs/swagger.json $(DOCKER_PATH)/public/docs/

docker.run:
	docker-compose -f docker/docker-compose.yml -p gfly up -d web

docker.logs:
	docker-compose -f docker/docker-compose.yml -p gfly logs -f web

docker.shell:
	docker-compose -f docker/docker-compose.yml -p gfly exec --user gfly web bash

docker.root:
	docker-compose -f docker/docker-compose.yml -p gfly exec web bash

docker.stop:
	docker-compose -f docker/docker-compose.yml -p gfly kill

docker.delete:
	docker-compose -f docker/docker-compose.yml -p gfly down

docker.destroy: docker.stop docker.delete

docker.build: docker.checking docker.vulncheck docker.test docker.shell

docker.release: docker.checking docker.test docker.shell
