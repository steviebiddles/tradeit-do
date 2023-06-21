package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"golang.org/x/exp/slog"
)

var httpAddr = ":3000"

func init() {
	if os.Getenv("HTTP_ADDR") != "" {
		httpAddr = os.Getenv("HTTP_ADDR")
	}
}

func main() {
	zl := zap.Must(zap.NewProduction())
	defer func(zl *zap.Logger) {
		_ = zl.Sync()
	}(zl)

	bi, _ := debug.ReadBuildInfo()
	sl := slog.New(zapslog.NewHandler(zl.Core())).With(
		slog.Group("program",
			slog.String("app", "tradeit"),
			slog.Int("pid", os.Getpid()),
			slog.String("go", bi.GoVersion),
		),
	)

	sl.Info("running tradeit")

	if err := run(sl); err != nil {
		sl.Error("exit tradeit", "error", err)
		os.Exit(1)
	}

	sl.Info("shutdown tradeit")
}

func run(sl *slog.Logger) error {
	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/live"))
	r.Use(middleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			sl.Info("homepage")

			w.Write([]byte("tradeit - homepage"))
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		sl.Info("page not found", "url", r.URL)

		w.WriteHeader(404)
		w.Write([]byte("tradeit - route not found"))
	})

	err := http.ListenAndServe(fmt.Sprintf("%s", httpAddr), r)
	if err != nil {
		return err
	}

	return nil
}
