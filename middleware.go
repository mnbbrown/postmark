package postmark

import (
	"github.com/mnbbrown/engine"
	"golang.org/x/net/context"
	"net/http"
)

type key int

const ctxKey int = 1

func Middleware(c *Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ctx := engine.GetContext(req).Set(ctxKey, c)
			next.ServeHTTP(rw, req)
		})
	}
}

func FromContext(ctx context.Context) (Client, bool) {
	cfg, ok := ctx.Value(ctxKey).(Client)
	return cfg, ok
}
