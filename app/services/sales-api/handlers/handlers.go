package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/HMadhav/service/app/services/sales-api/handlers/debug/checkgrp"
	"github.com/HMadhav/service/app/services/sales-api/handlers/v1/testgrp"
	"github.com/HMadhav/service/business/sys/auth"
	"github.com/HMadhav/service/business/web/mid"
	"github.com/HMadhav/service/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

func DebugMux(build string, log *zap.SugaredLogger, db *sqlx.DB) http.Handler {
	mux := DebugStandardLibraryMux()

	// Register debug check endpoints.
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
		DB:    db,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Auth     *auth.Auth
	Log      *zap.SugaredLogger
	DB       *sqlx.DB
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors((cfg.Log)), mid.Metrics(), mid.Panics())

	v1(app, cfg)

	return app
}

// Routes binds all the version 1 routes.
func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"

	tgh := testgrp.Handlers{
		Log: cfg.Log,
	}
	app.Handle(http.MethodGet, version, "/test", tgh.Test)
	app.Handle(http.MethodGet, version, "/testauth", tgh.Test, mid.Authenticate(cfg.Auth), mid.Authorize("ADMIN"))
}
