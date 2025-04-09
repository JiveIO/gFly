# gFly v1.12.1

**Laravel inspired web framework written in Go**

Built on top of [FastHttp - the fastest HTTP engine](https://github.com/valyala/fasthttp), [FluentSQL - flexible and powerful SQL builder](https://github.com/jiveio/fluentsql). Quick development with zero memory allocation and high performance. Very simple and easy to use.

# Tour of gFly

## I. Install environment

### 1. Install Docker [Docker Desktop](https://www.docker.com/products/docker-desktop/) or [OrbStack](https://orbstack.dev/)

### 2. Install Golang

### 2.1 On Mac
```bash
# Install go at folder /Users/$USER/Apps
mkdir -p /Users/$USER/Apps
wget https://go.dev/dl/go1.24.2.darwin-arm64.tar.gz
tar -xvzf go1.24.2.darwin-arm64.tar.gz
```
Add bottom of file `~/.profile` or `~/.zshrc`
```bash
vi ~/.profile

# ----------- Golang -----------
export GOROOT=/Users/$USER/Apps/go
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

### 2.2 On Linux
```bash
# Install go at folder /home/$USER/Apps
mkdir -p /home/$USER/Apps
wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz
tar -xvzf go1.24.2.linux-amd64.tar.gz
```
Add bottom of file `~/.profile` or `~/.zshrc`
```bash
vi ~/.profile

# ----------- Golang -----------
export GOROOT=/home/$USER/Apps/go
export GOPATH=/home/$USER/go
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

### 3. Install `Swag`, `Migrage`, `Lint`, and `Air`
```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Check Swag version
swag -v

# Install air
go install github.com/air-verse/air@latest

# Check Air version
air -v

# Install migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Check Migrate version
migrate --version

# Install Lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.0.2

# Check Lint version
golangci-lint --version
```

### 4. Create project skeleton from `gFly` repository
```bash
git clone https://github.com/jiveio/gfly.git myweb
cd myweb 
rm -rf .git* && cp .env.example .env
```

## II. Start `redis`, `mail`, and `db` services and `application`

Make sure don't have any services ran at ports `6379`, `1025`, `8025`, and `5432` on local. 

### 1. Start docker services
```bash
# Start
make docker.run
```
### 2. Check services
```bash
❯ docker ps

>>> CONTAINER ID   IMAGE                  COMMAND                  CREATED         STATUS                   PORTS                                                                                            NAMES
>>> 38fb5bd004df   redis:7.4.0-alpine     "docker-entrypoint.s…"   9 minutes ago   Up 9 minutes             0.0.0.0:6379->6379/tcp, :::6379->6379/tcp                                                        gfly-redis
>>> 9e52bdb5a4ae   axllent/mailpit        "/mailpit"               9 minutes ago   Up 9 minutes (healthy)   0.0.0.0:1025->1025/tcp, :::1025->1025/tcp, 0.0.0.0:8025->8025/tcp, :::8025->8025/tcp, 1110/tcp   gfly-mail
>>> d62e30b0d548   postgres:16.4-alpine   "docker-entrypoint.s…"   9 minutes ago   Up 9 minutes (healthy)   0.0.0.0:5432->5432/tcp, :::5432->5432/tcp                                                        gfly-db
```

### 3. Start app
```bash
# Doc
make doc

# Build (Build CLI and Web)
make build

# Docker run (Create DB, Redis, Mail services)
make docker.run

# Run
make dev
```

### 4. Check app

Browse URL [http://localhost:7789/](http://localhost:7789/)

Check API  via CLI
```
curl -v -X GET http://localhost:7789/api/v1/info | jq
```

Note: Install [jq](https://jqlang.github.io/jq/) tool to view JSON format

API doc [http://localhost:7789/docs/](http://localhost:7789/docs/)

### 5. CLI Actions

```bash
# Run command `hello-world`
./build/artisan cmd:run hello-world

# Run schedule 
./build/artisan schedule:run

# Run queue 
./build/artisan queue:run
```

Note: You can check more detail about [command](https://doc.gfly.dev/docs/03-digging-deeper/03-01-02.command/), [schedule](https://doc.gfly.dev/docs/03-digging-deeper/03-01-03.schedule/), and [queue](https://doc.gfly.dev/docs/03-digging-deeper/03-01-04.queue/) at link [https://doc.gfly.dev/](https://doc.gfly.dev/)


## III. Service connection

Add some code to check `application` connect to services `redis`, `mail`, and `db`.

### 1. Connect `Database` service

#### Import initial tables
```bash
make migrate.up
```

Note: Check DB connection and see 4 tables: `users`, `roles`, `user_roles`, and `address`.

#### Create command

Create a new command line `db-test`. Add file `app/console/commands/db_command.go` 

```go
package commands

import (
	"gfly/app/domain/models"
	"github.com/gflydev/console"
	"github.com/gflydev/core/log"
	mb "github.com/gflydev/db"
	"time"
)

// ---------------------------------------------------------------
//                      Register command.
// ./artisan cmd:run db-test
// ---------------------------------------------------------------

// Auto-register command.
func init() {
	console.RegisterCommand(&DBCommand{}, "db-test")
}

// ---------------------------------------------------------------
//                      DBCommand struct.
// ---------------------------------------------------------------

// DBCommand struct for hello command.
type DBCommand struct {
    console.Command
}

// Handle Process command.
func (c *DBCommand) Handle() {
    user, err := mb.GetModelBy[models.User]("email", "admin@gfly.dev")
    if err != nil || user == nil {
        log.Panic(err)
    }
    log.Infof("User %v\n", user)

    log.Infof("DBCommand :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}
```

#### Build and run command

```bash
# Build
make build

# Run command `db-test`
 ./build/artisan cmd:run db-test
```

### 2. Connect `Redis` service

Create a new command line `redis-test`. Add file `app/console/commands/redis_command.go` 

```go
package commands

import (
	"github.com/gflydev/cache"
	"github.com/gflydev/console"
	"github.com/gflydev/core/log"
	"time"
)

// ---------------------------------------------------------------
//                      Register command.
// ./artisan cmd:run redis-test
// ---------------------------------------------------------------

// Auto-register command.
func init() {
    console.RegisterCommand(&RedisCommand{}, "redis-test")
}

// ---------------------------------------------------------------
//                      RedisCommand struct.
// ---------------------------------------------------------------

// RedisCommand struct for hello command.
type RedisCommand struct {
    console.Command
}

// Handle Process command.
func (c *RedisCommand) Handle() {
    // Add new key
    if err := cache.Set("foo", "Hello world", time.Duration(15*24*3600)*time.Second); err != nil {
        log.Error(err)
    }

    // Get data key
    bar, err := cache.Get("foo")
    if err != nil {
        log.Error(err)
    }
    log.Infof("foo `%v`\n", bar)

    log.Infof("RedisCommand :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}

```

#### Build and run command

```bash
# Build
make build

# Run command `redis-test`
 ./build/artisan cmd:run redis-test
```

### 3. Connect `Mail` service

Create a new command line `mail-test`. Add file `app/console/commands/mail_command.go` 

```go
package commands

import (
	"github.com/gflydev/cache"
	"github.com/gflydev/console"
	"github.com/gflydev/core/log"
	"time"
)

// ---------------------------------------------------------------
//                      Register command.
// ./artisan cmd:run mail-test
// ---------------------------------------------------------------

// Auto-register command.
func init() {
    console.RegisterCommand(&RedisCommand{}, "mail-test")
}

// ---------------------------------------------------------------
//                      MailCommand struct.
// ---------------------------------------------------------------

// MailCommand struct for hello command.
type MailCommand struct {
    console.Command
}

// Handle Process command.
func (c *MailCommand) Handle() {
    // ============== Send mail ==============
    resetPassword := notifications.ResetPassword{
        Email: "admin@gfly.dev",
    }

    if err := notification.Send(resetPassword); err != nil {
        log.Error(err)
    }
  
    log.Infof("RedisCommand :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}
```

#### Build and run command

```bash
# Build
make build

# Run command `mail-test`
 ./build/artisan cmd:run mail-test
```

Check mail at http://localhost:8025/

### 4. Run command and check Log

```bash
curl -X 'GET' http://localhost:7789/api/v1/info
```

Check email at `http://localhost:8025/`
