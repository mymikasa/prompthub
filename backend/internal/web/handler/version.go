package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mymikasa/prompthub/internal/service"
	"github.com/mymikasa/prompthub/internal/web/result"
	"github.com/mymikasa/prompthub/pkg/ctxutil"
)

type VersionHandler struct {
	versionSvc *service.VersionService
}

func NewVersionHandler(versionSvc *service.VersionService) *VersionHandler {
	return &VersionHandler{versionSvc: versionSvc}
}

func (h *VersionHandler) Register(rg *gin.RouterGroup) {
	prompts := rg.Group("/api/prompts")
	prompts.GET("/:id/versions", h.List)
	prompts.GET("/:id/versions/:versionId", h.Get)
	prompts.POST("/:id/versions/:versionId/restore", h.Restore)
}

func (h *VersionHandler) List(c *gin.Context) {
	promptID, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid prompt id")
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	versions, err := h.versionSvc.ListVersions(c.Request.Context(), actor, promptID)
	if err != nil {
		if errors.Is(err, service.ErrPromptNotFound) {
			result.NotFound(c, "prompt not found")
			return
		}
		result.InternalError(c, err.Error())
		return
	}

	result.OK(c, versions)
}

func (h *VersionHandler) Get(c *gin.Context) {
	promptID, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid prompt id")
		return
	}

	versionID, err := strconv.ParseInt(c.Param("versionId"), 10, 64)
	if err != nil {
		result.BadRequest(c, "invalid version id")
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	v, err := h.versionSvc.GetVersion(c.Request.Context(), actor, promptID, versionID)
	if err != nil {
		if errors.Is(err, service.ErrVersionNotFound) || errors.Is(err, service.ErrPromptNotFound) {
			result.NotFound(c, "version not found")
			return
		}
		result.InternalError(c, err.Error())
		return
	}

	result.OK(c, v)
}

func (h *VersionHandler) Restore(c *gin.Context) {
	promptID, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid prompt id")
		return
	}

	versionID, err := strconv.ParseInt(c.Param("versionId"), 10, 64)
	if err != nil {
		result.BadRequest(c, "invalid version id")
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	v, err := h.versionSvc.RestoreVersion(c.Request.Context(), actor, promptID, versionID)
	if err != nil {
		if errors.Is(err, service.ErrVersionNotFound) || errors.Is(err, service.ErrPromptNotFound) {
			result.NotFound(c, "version not found")
			return
		}
		if errors.Is(err, service.ErrNoPermission) {
			result.Forbidden(c, err.Error())
			return
		}
		result.InternalError(c, err.Error())
		return
	}

	result.OK(c, v)
}
