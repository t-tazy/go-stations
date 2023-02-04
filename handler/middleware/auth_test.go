package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TechBowl-japan/go-stations/config"
)

func TestNewAuthMiddleware(t *testing.T) {
	// 環境変数のダミー
	cfg := &config.Config{
		UserID:   "1111",
		Password: "2222",
	}

	mux := http.NewServeMux()
	authmiddleware := NewAuthMiddleware(cfg)

	health := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		rsp := struct {
			Message string
		}{"health"}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			t.Fatalf("%v", err)
		}
	})
	mux.Handle("/health", health)

	todo := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		rsp := struct {
			Message string
		}{"todo"}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			t.Fatalf("%v", err)
		}
	})
	mux.Handle("/todos", authmiddleware(todo))

	type req struct {
		path     string
		auth     bool
		userID   string
		password string
	}
	tests := map[string]struct {
		req        req
		wantStatus int
	}{
		"ok": {
			req: req{
				path:     "/todos",
				auth:     true,
				userID:   cfg.UserID,
				password: cfg.Password,
			},
			wantStatus: http.StatusOK,
		},
		"request public endpoint": {
			req: req{
				path: "/health",
				auth: false,
			},
			wantStatus: http.StatusOK,
		},
		"wrong user_id, password": {
			req: req{
				path:     "/todos",
				auth:     true,
				userID:   "3333",
				password: "4444",
			},
			wantStatus: http.StatusUnauthorized,
		},
		"empty user_id, password": {
			req: req{
				path: "/todos",
				auth: true,
			},
			wantStatus: http.StatusUnauthorized,
		},
		"don't include Authorization header": {
			req: req{
				path: "/todos",
				auth: false,
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for key, test := range tests {
		test := test
		t.Run(key, func(t *testing.T) {

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodGet,
				test.req.path,
				nil,
			)

			if test.req.auth {
				r.SetBasicAuth(test.req.userID, test.req.password)
			}
			mux.ServeHTTP(w, r)
			rsp := w.Result()
			if rsp.StatusCode != test.wantStatus {
				t.Errorf("want status %d, but got %d", test.wantStatus, rsp.StatusCode)
			}
		})
	}
}
