package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mymikasa/prompthub/internal/config"
	"github.com/mymikasa/prompthub/internal/repo"
	"github.com/mymikasa/prompthub/internal/repo/dao"
	"github.com/mymikasa/prompthub/internal/service"
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

	if err := ioc.Migrate(db); err != nil {
		slog.Error("run migrate", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if cfg.IsDev() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	userDAO := dao.NewUserDAO(db)
	workspaceDAO := dao.NewWorkspaceDAO(db)
	promptDAO := dao.NewPromptDAO(db)
	tagDAO := dao.NewTagDAO(db)
	variableDAO := dao.NewVariableDAO(db)

	userRepo := repo.NewUserRepo(userDAO)
	workspaceRepo := repo.NewWorkspaceRepo(workspaceDAO)
	promptRepo := repo.NewPromptRepo(promptDAO)
	tagRepo := repo.NewTagRepo(tagDAO)
	variableRepo := repo.NewVariableRepo(variableDAO)

	authService := service.NewAuthService(userRepo, workspaceRepo)
	promptService := service.NewPromptService(promptRepo, tagRepo, variableRepo)
	tagService := service.NewTagService(tagRepo)
	variableService := service.NewVariableService(variableRepo)

	healthHandler := handler.NewHealthHandler(db)
	authHandler := handler.NewAuthHandler(authService, cfg.SessionSecret, cfg.IsDev())
	promptHandler := handler.NewPromptHandler(promptService)
	tagHandler := handler.NewTagHandler(tagService)
	variableHandler := handler.NewVariableHandler(variableService, promptService)

	engine := gin.New()
	router.Setup(engine, healthHandler, authHandler, promptHandler, tagHandler, variableHandler, cfg.SessionSecret, userRepo, workspaceRepo)

	slog.Info("server starting", slog.String("addr", cfg.HTTPAddr))
	if err := engine.Run(cfg.HTTPAddr); err != nil {
		slog.Error("server stopped", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
