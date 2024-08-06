package api

import (
	"app/app/http/response"
	"app/core/gfly"
	"fmt"
	"os"
)

// ===========================================================================================================
// 											Default API
// ===========================================================================================================

// NewDefaultApi As a constructor to create new API.
func NewDefaultApi() *DefaultApi {
	return &DefaultApi{}
}

// DefaultApi API struct.
type DefaultApi struct {
	gfly.Api
}

// Handle Process main logic for API.
// @Summary Get API info
// @Description Get API server information
// @Tags Misc
// @Accept json
// @Produce json
// @Success 200 {object} response.ServerInfo
// @Router /info [get]
func (h *DefaultApi) Handle(c *gfly.Ctx) error {
	obj := response.ServerInfo{
		Name: os.Getenv("API_NAME"),
		Prefix: fmt.Sprintf(
			"/%s/%s",
			os.Getenv("API_PREFIX"),
			os.Getenv("API_VERSION"),
		),
		Server: os.Getenv("APP_NAME"),
	}

	return c.JSON(obj)
}
