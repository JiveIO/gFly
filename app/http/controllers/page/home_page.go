package page

import "github.com/gflydev/core"

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

// NewHomePage As a constructor to create a Home Page.
func NewHomePage() *HomePage {
	return &HomePage{}
}

type HomePage struct {
	core.Page
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

func (m *HomePage) Handle(c *core.Ctx) error {
	return c.Redirect("/index.html")
}
