package routes

import (
	"app/core/gfly"
)

// AppRoutes func for describe all routes in system.
func AppRoutes(f gfly.IFly) {
	ApiRoutes(f) // Register API routes.
	WebRoutes(f) // Register Web routes.
}
