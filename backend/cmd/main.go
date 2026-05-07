package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mymikasa/prompthub/internal/config"
	"github.com/mymikasa/prompthub/internal/web/handler"
	"github.com/mymikasa/prompthub/internal/web/router"
	"github.com/mymikasa/prompthub/ioc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	db, err := ioc.InitDB(cfg)
	if err != nil {
		slog.Error("init db", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if cfg.IsDev() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	healthHandler := handler.NewHealthHandler(db)
	router.Setup(engine, healthHandler)

	slog.Info("server starting", slog.String("addr", cfg.HTTPAddr))
	if err := engine.Run(cfg.HTTPAddr); err != nil {
		slog.Error("server stopped", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
