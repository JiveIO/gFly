package page

import (
	"app/core/gfly"
)

// ===========================================================================================================
// 											Home page
// ===========================================================================================================

// NewHomePage As a constructor to create new Page.
func NewHomePage() *HomePage {
	return &HomePage{}
}

type HomePage struct {
	gfly.Page
}

func (m *HomePage) Handle(c *gfly.Ctx) error {
	return c.Redirect("/index.html")
}
