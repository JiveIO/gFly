package routes

import (
	"github.com/gflydev/core"
)

func Router(r core.IFlyRouter) {
	ApiRoutes(r) // Register API routes.
	WebRoutes(r) // Register Web routes.
}
