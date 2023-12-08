package main

import (
	"app/app/http/routes"
	"app/core/gfly"
	"app/core/middleware"
	"app/docs"
)

// Main function
// @title API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email vinh@jivecode.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server gFly."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "gFly.dev"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	app := gfly.New()

	// Add global middlewares
	app.Use(middleware.CORS(map[string]string{
		gfly.HeaderAccessControlAllowOrigin: "*",
	}))

	routes.AppRoutes(app)

	// curl -X GET http://localhost:7789/index.html
	app.ServeFiles()

	app.Start()
}
