package routes

import (
	"app/app/http/controllers/page"
	"app/core/gfly"
)

// WebRoutes func for describe a group of Web page routes.
func WebRoutes(f gfly.IFly) {
	// Web Routers
	f.GET("/", page.NewHomePage())
}
