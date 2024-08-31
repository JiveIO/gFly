package api

import (
	"fmt"
	"gfly/app/domain/repository"
	"gfly/app/http/response"
	"gfly/app/notifications"
	"github.com/gflydev/cache"
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/notification"
	"time"
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
	core.Api
}

// Handle Process main logic for API.
// @Summary Get API info
// @Description Get API server information
// @Tags Misc
// @Accept json
// @Produce json
// @Success 200 {object} response.ServerInfo
// @Router /info [get]
func (h *DefaultApi) Handle(c *core.Ctx) error {
	// ============== Check database ==============
	user := repository.Pool.GetUserByEmail("admin@gfly.dev")
	log.Infof("User %v\n", user)

	// ============== Set/Get caching ==============
	if err := cache.Set("foo", "Hello world", time.Duration(15*24*3600)*time.Second); err != nil {
		log.Error(err)
	}
	bar, err := cache.Get("foo")
	if err != nil {
		log.Error(err)
	}
	log.Infof("foo `%v`\n", bar)

	// ============== Send mail ==============
	resetPassword := notifications.ResetPassword{
		Email: "vinh@jivecode.com",
	}
	if err = notification.Send(resetPassword); err != nil {
		log.Error(err)
	}

	obj := response.ServerInfo{
		Name: utils.Getenv("API_NAME", "gfly"),
		Prefix: fmt.Sprintf(
			"/%s/%s",
			utils.Getenv("API_PREFIX", "api"),
			utils.Getenv("API_VERSION", "v1"),
		),
		Server: utils.Getenv("APP_URL", "http://localhost:7789"),
	}

	return c.JSON(obj)
}
