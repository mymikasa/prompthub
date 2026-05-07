package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mymikasa/prompthub/internal/repo"
	"github.com/mymikasa/prompthub/internal/web/handler"
	"github.com/mymikasa/prompthub/internal/web/middleware"
)

func Setup(engine *gin.Engine, health *handler.HealthHandler, auth *handler.AuthHandler, sessionSecret string, userRepo *repo.UserRepo) {
	engine.Use(middleware.Recovery())
	engine.Use(middleware.RequestLog())
	engine.Use(middleware.CORS())

	public := engine.Group("")
	authenticated := engine.Group("")
	authenticated.Use(middleware.Auth(sessionSecret, userRepo))

	health.Register(public)
	auth.Register(public, authenticated)
}
