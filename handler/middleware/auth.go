package middleware

import (
	"net/http"

	"github.com/TechBowl-japan/go-stations/config"
)

// Basic認証を行うmiddlewareを返す
func NewAuthMiddleware(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			userID, password, ok := r.BasicAuth()
			if !ok || cfg.UserID != userID || cfg.Password != password {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
