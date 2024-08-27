package page

import "github.com/gflydev/core"

// ===========================================================================================================
// 											Home page
// ===========================================================================================================

// NewHomePage As a constructor to create a Home Page.
func NewHomePage() *HomePage {
	return &HomePage{}
}

type HomePage struct {
	core.Page
}

func (m *HomePage) Handle(c *core.Ctx) error {
	return c.Redirect("/index.html")
}
