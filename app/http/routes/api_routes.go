package routes

import (
	"fmt"
	"gfly/app/http/controllers/api"
	"github.com/gflydev/core"
	"github.com/gflydev/core/utils"
)

// ApiRoutes func for describe a group of API routes.
func ApiRoutes(r core.IFlyRouter) {
	prefixAPI := fmt.Sprintf(
		"/%s/%s",
		utils.Getenv("API_PREFIX", "api"),
		utils.Getenv("API_VERSION", "v1"),
	)

	// API Routers
	r.Group(prefixAPI, func(r *core.Group) {
		// curl -v -X GET http://localhost:7789/api/v1/info | jq
		r.GET("/info", api.NewDefaultApi())
	})
}
