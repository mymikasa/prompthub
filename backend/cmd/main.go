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
	versionDAO := dao.NewVersionDAO(db)
	testCaseDAO := dao.NewTestCaseDAO(db)
	providerConfigDAO := dao.NewProviderConfigDAO(db)
	promptRunDAO := dao.NewPromptRunDAO(db)

	userRepo := repo.NewUserRepo(userDAO)
	workspaceRepo := repo.NewWorkspaceRepo(workspaceDAO)
	promptRepo := repo.NewPromptRepo(promptDAO)
	tagRepo := repo.NewTagRepo(tagDAO)
	variableRepo := repo.NewVariableRepo(variableDAO)
	versionRepo := repo.NewVersionRepo(versionDAO)
	testCaseRepo := repo.NewTestCaseRepo(testCaseDAO)
	providerConfigRepo := repo.NewProviderConfigRepo(providerConfigDAO)
	promptRunRepo := repo.NewPromptRunRepo(promptRunDAO)

	authService := service.NewAuthService(userRepo, workspaceRepo)
	promptService := service.NewPromptService(promptRepo, tagRepo, variableRepo, versionRepo, providerConfigRepo)
	tagService := service.NewTagService(tagRepo)
	variableService := service.NewVariableService(variableRepo)
	versionService := service.NewVersionService(versionRepo, promptRepo, tagRepo, variableRepo)
	testCaseService := service.NewTestCaseService(testCaseRepo, promptRepo, variableRepo)
	encryptionKey := cfg.SecretEncryptionKey
	if cfg.IsDev() {
		encryptionKey = "0000000000000000000000000000000000000000000000000000000000000000"
	}
	providerConfigService, err := service.NewProviderConfigService(providerConfigRepo, encryptionKey)
	if err != nil {
		slog.Error("init provider config service", slog.String("error", err.Error()))
		os.Exit(1)
	}
	promptRunService := service.NewPromptRunService(promptRunRepo, promptRepo, testCaseRepo, variableRepo, providerConfigService, promptService)

	healthHandler := handler.NewHealthHandler(db)
	authHandler := handler.NewAuthHandler(authService, cfg.SessionSecret, cfg.IsDev())
	promptHandler := handler.NewPromptHandler(promptService)
	tagHandler := handler.NewTagHandler(tagService)
	variableHandler := handler.NewVariableHandler(variableService, promptService)
	versionHandler := handler.NewVersionHandler(versionService)
	testCaseHandler := handler.NewTestCaseHandler(testCaseService)
	providerConfigHandler := handler.NewProviderConfigHandler(providerConfigService)
	promptRunHandler := handler.NewPromptRunHandler(promptRunService)

	engine := gin.New()
	router.Setup(engine, healthHandler, authHandler, promptHandler, tagHandler, variableHandler, versionHandler, testCaseHandler, providerConfigHandler, promptRunHandler, cfg.SessionSecret, userRepo, workspaceRepo)

	slog.Info("server starting", slog.String("addr", cfg.HTTPAddr))
	if err := engine.Run(cfg.HTTPAddr); err != nil {
		slog.Error("server stopped", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
