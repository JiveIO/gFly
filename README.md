# gFly v1.6.0

**Laravel inspired web framework written in Go**

Built on top of [FastHttp - the fastest HTTP engine](https://github.com/valyala/fasthttp), [FluentSQL - flexible and powerful SQL builder](https://github.com/jiveio/fluentsql). Quick development with zero memory allocation and high performance. Very simple and easy to use.

# Tour of gFly

## I. Install environment

### 1. Install Docker [Docker Desktop](https://www.docker.com/products/docker-desktop/) or [OrbStack](https://orbstack.dev/)

### 2. Install Golang
```bash
# Install go at folder /Users/$USER/Apps
mkdir -p /Users/$USER/Apps
wget https://go.dev/dl/go1.22.6.darwin-arm64.tar.gz
tar -xvzf go1.22.6.darwin-arm64.tar.gz

# So Go root path is /Users/$USER/Apps/go1.22.6
```
Add bottom of file `~/.profile` or `~/.zshrc`
```bash
vi ~/.profile

# ----------- Golang -----------
export GOROOT=/Users/$USER/Apps/go1.22.6
export GOPATH=/Users/$USER/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

Check
```bash
source ~/.profile
# Or
source ~/.zshrc

# Check Go
go version
```

### 3. Install `Swag` and `Air`
```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Check Swag version
swag -v

# Install air
go install github.com/air-verse/air@latest

# Check Air version
air -v
```

### 4. Create project skeleton from `gFly` repository
```bash
git clone https://github.com/jiveio/gfly.git
cd gFly
rm -rf .git*
cp .env.example .env
```

## II. Start `redis`, `mail`, and `db` services

Make sure don't have any services ran at ports `6379`, `1025`, `8025`, and `5432` on local. 

### 1. Start docker services
```bash
# Start
make docker.run
```
### 2. Check services
```bash
❯ docker ps

CONTAINER ID   IMAGE                  COMMAND                  CREATED         STATUS                   PORTS                                                                                            NAMES
38fb5bd004df   redis:7.4.0-alpine     "docker-entrypoint.s…"   9 minutes ago   Up 9 minutes             0.0.0.0:6379->6379/tcp, :::6379->6379/tcp                                                        gfly-redis
9e52bdb5a4ae   axllent/mailpit        "/mailpit"               9 minutes ago   Up 9 minutes (healthy)   0.0.0.0:1025->1025/tcp, :::1025->1025/tcp, 0.0.0.0:8025->8025/tcp, :::8025->8025/tcp, 1110/tcp   gfly-mail
d62e30b0d548   postgres:16.4-alpine   "docker-entrypoint.s…"   9 minutes ago   Up 9 minutes (healthy)   0.0.0.0:5432->5432/tcp, :::5432->5432/tcp                                                        gfly-db
```

### Start app
```bash
# Doc
make doc

# Run
go run main.go

# Air
air main.go
```

### Check app

Browse URL [http://localhost:7789/](http://localhost:7789/)

Check API `curl -v -X GET http://localhost:7789/api/v1/info | jq`

API doc [http://localhost:7789/docs/](http://localhost:7789/docs/)

Note: Install [jq](https://jqlang.github.io/jq/) tool to view JSON format
