package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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

type AccessLog struct {
	Timestamp time.Time `json:"timestamp"`
	Latency   int64     `json:"latency"` // milliseconds
	Path      string    `json:"path"`
	OS        string    `json:"os"`
}

// OsMiddleware適用後に適用すること
func AccessLogMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Now().Sub(start).Milliseconds()
		os, ok := GetOS(r.Context())
		if !ok {
			fmt.Println("OS not found")
		}
		jsonBytes, err := json.Marshal(AccessLog{
			Timestamp: start,
			Latency:   duration,
			Path:      r.URL.Path,
			OS:        os,
		})
		if err != nil {
			log.Printf("%v", err)
			return
		}
		fmt.Println(string(jsonBytes))
	}
	return http.HandlerFunc(fn)
}
