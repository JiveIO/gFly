package routes

import (
	"gfly/app/http/controllers/page"
	"github.com/gflydev/core"
)

// WebRoutes func for describe a group of Web page routes.
func WebRoutes(r core.IFlyRouter) {
	// Web Routers
	r.GET("/", page.NewHomePage())
}
