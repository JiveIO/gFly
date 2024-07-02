/*
Package middleware provides set of middleware functions that can be used to authenticate and authorize
requests in HTTP server.It also supports handling CORS, propagating headers, integrating with New Relic APM, and enabling
distributed tracing using OpenTelemetry.
*/
package middleware

import (
	"app/core/gfly"
)

const (
	AllowedOrigin  = "*"
	AllowedHeaders = "Authorization, Content-Type, x-requested-with, origin, true-client-ip, X-Correlation-ID"
	AllowedMethods = "PUT, POST, GET, DELETE, OPTIONS, PATCH"
)

// CORS an HTTP middleware that sets headers based on the provided envHeaders configuration
//
// Example: Add global middlewares in main file
//
//	app.Use(middleware.CORS(map[string]string{
//		gfly.HeaderAccessControlAllowOrigin: "*",
//	}))
func CORS(envHeaders map[string]string) gfly.MiddlewareHandler {
	return func(c *gfly.Ctx) error {
		corsHeadersConfig := getValidCORSHeaders(envHeaders)
		for k, v := range corsHeadersConfig {
			c.SetHeader(k, v)
		}

		return nil
	}
}

// getValidCORSHeaders returns a validated map of CORS headers.
// values specified in env are present in envHeaders
func getValidCORSHeaders(envHeaders map[string]string) map[string]string {
	validCORSHeadersAndValues := make(map[string]string)

	for _, header := range allowedCORSHeader() {
		// If config is set, use that
		if val, ok := envHeaders[header]; ok && val != "" {
			validCORSHeadersAndValues[header] = val
			continue
		}

		// If config is not set - for the three headers, set default value.
		switch header {
		case gfly.HeaderAccessControlAllowOrigin:
			validCORSHeadersAndValues[header] = AllowedOrigin
		case gfly.HeaderAccessControlAllowHeaders:
			validCORSHeadersAndValues[header] = AllowedHeaders
		case gfly.HeaderAccessControlAllowMethods:
			validCORSHeadersAndValues[header] = AllowedMethods
		}
	}

	val := validCORSHeadersAndValues[gfly.HeaderAccessControlAllowHeaders]

	if val != AllowedHeaders {
		validCORSHeadersAndValues[gfly.HeaderAccessControlAllowHeaders] = AllowedHeaders + ", " + val
	}

	return validCORSHeadersAndValues
}

// allowedCORSHeader returns the HTTP headers used for CORS configuration in web applications.
func allowedCORSHeader() []string {
	return []string{
		gfly.HeaderAccessControlAllowOrigin,
		gfly.HeaderAccessControlAllowHeaders,
		gfly.HeaderAccessControlAllowMethods,
		gfly.HeaderAccessControlAllowCredentials,
		gfly.HeaderAccessControlExposeHeaders,
		gfly.HeaderAccessControlMaxAge,
	}
}
