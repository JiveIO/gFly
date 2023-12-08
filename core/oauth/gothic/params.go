package gothic

import (
	"app/core/utils"
	"github.com/valyala/fasthttp"
)

type Params struct {
	ctx *fasthttp.RequestCtx
}

func (p *Params) Get(key string) string {
	return utils.UnsafeString(p.ctx.QueryArgs().Peek(key))
}
