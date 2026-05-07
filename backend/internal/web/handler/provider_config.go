package handler

import (
	"errors"
	"strconv"

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
	settings.GET("/providers", h.List)
	settings.POST("/providers", h.Save)
	settings.GET("/providers/:id", h.GetByID)
	settings.DELETE("/providers/:id", h.Delete)
}

func (h *ProviderConfigHandler) List(c *gin.Context) {
	actor := ctxutil.ActorFromCtx(c.Request.Context())
	resp, err := h.svc.List(c.Request.Context(), actor)
	if err != nil {
		result.InternalError(c, err.Error())
		return
	}
	result.OK(c, resp)
}

func (h *ProviderConfigHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		result.BadRequest(c, "invalid id")
		return
	}
	actor := ctxutil.ActorFromCtx(c.Request.Context())
	resp, err := h.svc.GetByID(c.Request.Context(), actor, id)
	if err != nil {
		if errors.Is(err, service.ErrProviderConfigNotFound) {
			result.NotFound(c, "provider config not found")
			return
		}
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

func (h *ProviderConfigHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		result.BadRequest(c, "invalid id")
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	if err := h.svc.Delete(c.Request.Context(), actor, id); err != nil {
		if errors.Is(err, service.ErrProviderConfigNotFound) {
			result.NotFound(c, "provider config not found")
			return
		}
		result.InternalError(c, err.Error())
		return
	}
	result.OK(c, nil)
}
