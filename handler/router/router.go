package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	h := handler.NewHealthzHandler()
	mux.Handle("/healthz", middleware.OsMiddleware(middleware.AccessLogMiddleware(h)))

	svc := service.NewTODOService(todoDB)
	th := handler.NewTODOHandler(svc)
	mux.Handle("/todos", th)

	p := handler.NewPanicHandler()
	// mxu.Handle("/do-panic", p)
	mux.Handle("/do-panic", middleware.RecoveryMiddleware(p))

	return mux
}
