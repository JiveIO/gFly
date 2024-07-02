package api

import (
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
// @Description Get API information
// @Tags Mics
// @Accept json
// @Produce json
// @Success 200
// @Router /info [get]
func (h *DefaultApi) Handle(c *gfly.Ctx) error {
	obj := map[string]any{
		"name": os.Getenv("API_NAME"),
		"prefix": fmt.Sprintf(
			"/%s/%s",
			os.Getenv("API_PREFIX"),
			os.Getenv("API_VERSION"),
		),
		"server": os.Getenv("APP_NAME"),
	}

	return c.JSON(obj)
}
