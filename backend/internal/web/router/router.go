package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mymikasa/prompthub/internal/web/handler"
	"github.com/mymikasa/prompthub/internal/web/middleware"
)

func Setup(engine *gin.Engine, health *handler.HealthHandler) {
	engine.Use(middleware.Recovery())
	engine.Use(middleware.RequestLog())

	api := engine.Group("")
	health.Register(api)
}
