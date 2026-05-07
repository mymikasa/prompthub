package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mymikasa/prompthub/internal/service"
	"github.com/mymikasa/prompthub/internal/web/result"
	"github.com/mymikasa/prompthub/pkg/ctxutil"
)

type ProviderConfigHandler struct {
	svc *service.ProviderConfigService
}

func NewProviderConfigHandler(svc *service.ProviderConfigService) *ProviderConfigHandler {
	return &ProviderConfigHandler{svc: svc}
}

func (h *ProviderConfigHandler) Register(rg *gin.RouterGroup) {
	settings := rg.Group("/api/settings")
	settings.GET("/provider", h.Get)
	settings.PUT("/provider", h.Save)
}

func (h *ProviderConfigHandler) Get(c *gin.Context) {
	actor := ctxutil.ActorFromCtx(c.Request.Context())
	resp, err := h.svc.Get(c.Request.Context(), actor)
	if err != nil {
		result.InternalError(c, err.Error())
		return
	}
	result.OK(c, resp)
}

func (h *ProviderConfigHandler) Save(c *gin.Context) {
	var req service.SaveProviderConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, err.Error())
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	resp, err := h.svc.Save(c.Request.Context(), actor, &req)
	if err != nil {
		result.InternalError(c, err.Error())
		return
	}

	result.OK(c, resp)
}
