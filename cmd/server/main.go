package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/soulteary/gorge-highlight/internal/config"
	"github.com/soulteary/gorge-highlight/internal/highlight"
	"github.com/soulteary/gorge-highlight/internal/httpapi"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var cfg *config.Config

	if path := os.Getenv("HIGHLIGHT_CONFIG_FILE"); path != "" {
		var err error
		cfg, err = config.LoadFromFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load config from %s: %v\n", path, err)
			os.Exit(1)
		}
	} else {
		cfg = config.LoadFromEnv()
	}

	hl := highlight.New()

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true, LogURI: true, LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			slog.Info("REQUEST", "method", v.Method, "uri", v.URI, "status", v.Status)
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("2M"))

	httpapi.RegisterRoutes(e, &httpapi.Deps{
		Highlighter: hl,
		Token:       cfg.ServiceToken,
		MaxBytes:    cfg.MaxBytes,
	})

	e.Logger.Fatal(e.Start(cfg.ListenAddr))
}
