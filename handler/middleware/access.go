package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

type agentOSKey struct{} // contextのkey

// contextからOS情報を取得する
func GetOS(ctx context.Context) (string, bool) {
	os, ok := ctx.Value(agentOSKey{}).(string)
	return os, ok
}

func OsMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ua := useragent.Parse(r.UserAgent())
		ctx = context.WithValue(ctx, agentOSKey{}, ua.OS)
		clone := r.Clone(ctx)
		next.ServeHTTP(w, clone)
	}
	return http.HandlerFunc(fn)
}
