package postmark

import (
	"github.com/mnbbrown/engine"
	"net/http"
)

type key int

const ctxKey key = 1

func Middleware(c *Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			engine.GetContext(req).Set(ctxKey, c)
			next.ServeHTTP(rw, req)
		})
	}
}

func FromContext(ctx *engine.Context) (*Client, bool) {
	cfg, ok := ctx.Value(ctxKey).(*Client)
	return cfg, ok
}
