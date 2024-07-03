# gFly v1.5.0

**Laravel inspired web framework written in Go**

Built on top of [FastHttp - the fastest HTTP engine](https://github.com/valyala/fasthttp), [FluentSQL - flexible and powerful SQL builder](https://github.com/jiveio/fluentsql) and [FluentModel - flexible and powerful Data-Access Layer](https://github.com/jiveio/fluentmodel) for Go. Quick development with zero memory allocation and high performance. Very simple and easy to use.

# Tour of gFly

### Install environment

Assume you finish step [Installation](https://doc.gfly.dev/docs/01-greeting-start/01-01-01.installation/) by using [Docker Setup](https://doc.gfly.dev/docs/01-greeting-start/01-01-01.installation/#i-docker-setup) or [Local Setup](https://doc.gfly.dev/docs/01-greeting-start/01-01-01.installation/#ii-local-setup)

First of all. Need to start your application by below commands:

### Pull `gFly` code structure
```bash
git clone https://github.com/JiveIO/gFly.git && cd gFly && rm -rf .git* && cp .env.example .env
```

### Start app
```bash
# Local setup
go run main.go
```

### Check app

Browse URL [http://localhost:7789/](http://localhost:7789/)

Check API `curl -v -X GET http://localhost:7789/api/v1/info | jq`

API doc [http://localhost:7789/docs/](http://localhost:7789/docs/)

# Create Page and API

We will experiment with the example of creating a page `HelloPage` and an API `HelloApi` controllers.

## Create a HelloPage

### Controller

Create file `hello_page.go` in folder `app/http/controllers/page/`

```go
package page

import (
    "app/core/gfly"
)

// ===============================================================================
//                                  Hello page
// ===============================================================================

// NewHelloPage As a constructor to create new Page.
func NewHelloPage() *HelloPage {
    return &HelloPage{}
}

type HelloPage struct {
    gfly.Page
}

func (m *HelloPage) Handle(c *gfly.Ctx) error {
    return c.HTML("<h2>Hello world</h2>")
}
```

### Define router

Declare `HelloPage` into `app/http/routes/web_routes.go`

```go
package routes

import (
    "app/app/http/controllers/page"
    "app/core/gfly"
)

// WebRoutes func for describe a group of Web page routes.
func WebRoutes(f gfly.IFly) {
    // Web Routers
    f.GET("/hello", page.NewHelloPage())
}

```

### Checking

Browse [http://localhost:7789/hello](http://localhost:7789/hello)

# Create a HelloAPI


### API

Create file `hello_api.go` in folder `app/http/controllers/api/`

```go
package api

import (
    "app/core/gfly"
    "fmt"
    "os"
)

// ===============================================================================
//                                  Hello API
// ===============================================================================


// NewHelloApi As a constructor to create new API.
func NewHelloApi() *HelloApi {
    return &HelloApi{}
}

// HelloApi API struct.
type HelloApi struct {
    gfly.Api
}

// Handle Process main logic for API.
func (h *HelloApi) Handle(c *gfly.Ctx) error {
    obj := map[string]any{
        "name": os.Getenv("API_NAME"),
        "server": os.Getenv("APP_NAME"),
    }

    return c.JSON(obj)
}
```

### Define router

Declare `HelloApi` into `app/http/routes/api_routes.go`

```go
package routes

import (
    "app/app/http/controllers/api"
    "app/core/gfly"
    "fmt"
    "os"
)

// ApiRoutes func for describe a group of API routes.
func ApiRoutes(f gfly.IFly) {
    prefixAPI := fmt.Sprintf(
        "/%s/%s",
        os.Getenv("API_PREFIX"),
        os.Getenv("API_VERSION"),
    )

    f.Group(prefixAPI, func(apiRouter *gfly.Group) {
        // curl -v -X GET http://localhost:7789/api/v1/hello | jq
        apiRouter.GET("/hello", api.NewHelloApi())
    })
}
```

### Checking

Browse [http://localhost:7789/api/v1/hello](http://localhost:7789/api/v1/hello)

Or terminal
```bash
curl -v -X GET http://localhost:7789/api/v1/hello | jq
```

Note: Install [jq](https://jqlang.github.io/jq/) tool to view JSON format