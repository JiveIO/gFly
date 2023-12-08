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
		// API Routers
		// curl -v -X GET http://localhost:7789/api/v1/info | jq
		apiRouter.GET("/info", api.NewDefaultApi())
	})
}
