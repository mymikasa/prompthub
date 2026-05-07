package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/service"
	"github.com/mymikasa/prompthub/internal/web/result"
	"github.com/mymikasa/prompthub/pkg/ctxutil"
)

type PromptHandler struct {
	promptSvc *service.PromptService
}

func NewPromptHandler(promptSvc *service.PromptService) *PromptHandler {
	return &PromptHandler{promptSvc: promptSvc}
}

func (h *PromptHandler) Register(rg *gin.RouterGroup) {
	prompts := rg.Group("/api/prompts")
	prompts.POST("", h.Create)
	prompts.GET("", h.List)
	prompts.GET("/:id", h.Get)
	prompts.PATCH("/:id", h.Update)
	prompts.POST("/:id/archive", h.Archive)
	prompts.POST("/:id/restore", h.Restore)
	prompts.POST("/:id/render", h.Render)
}

func (h *PromptHandler) Create(c *gin.Context) {
	var req service.CreatePromptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, err.Error())
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	p, err := h.promptSvc.CreatePrompt(c.Request.Context(), actor, &req)
	if err != nil {
		result.InternalError(c, err.Error())
		return
	}

	result.Created(c, p)
}

func (h *PromptHandler) List(c *gin.Context) {
	actor := ctxutil.ActorFromCtx(c.Request.Context())
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "pageSize", 20)

	filter := service.PromptFilter{
		Keyword:  c.Query("keyword"),
		Provider: c.Query("provider"),
		Model:    c.Query("model"),
	}
	if statuses := c.QueryArray("status"); len(statuses) > 0 {
		filter.Statuses = statuses
	}
	if tags := c.QueryArray("tags"); len(tags) > 0 {
		filter.Tags = tags
	}

	list, err := h.promptSvc.ListPrompts(c.Request.Context(), actor, page, pageSize, filter)
	if err != nil {
		result.InternalError(c, err.Error())
		return
	}

	result.OK(c, gin.H{
		"items":    list.Items,
		"total":    list.Total,
		"page":     list.Page,
		"pageSize": list.PageSize,
	})
}

func (h *PromptHandler) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid id")
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	p, err := h.promptSvc.GetPrompt(c.Request.Context(), actor, id)
	if err != nil {
		if errors.Is(err, service.ErrPromptNotFound) {
			result.NotFound(c, "prompt not found")
			return
		}
		result.InternalError(c, err.Error())
		return
	}

	result.OK(c, p)
}

func (h *PromptHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid id")
		return
	}

	var req service.UpdatePromptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, err.Error())
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	p, err := h.promptSvc.UpdatePrompt(c.Request.Context(), actor, id, &req)
	if err != nil {
		writePromptError(c, err)
		return
	}

	result.OK(c, p)
}

func (h *PromptHandler) Archive(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid id")
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	if err := h.promptSvc.ArchivePrompt(c.Request.Context(), actor, id); err != nil {
		writePromptError(c, err)
		return
	}

	result.OK(c, nil)
}

func (h *PromptHandler) Restore(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid id")
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	if err := h.promptSvc.RestorePrompt(c.Request.Context(), actor, id); err != nil {
		writePromptError(c, err)
		return
	}

	result.OK(c, nil)
}

func (h *PromptHandler) Render(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid id")
		return
	}

	var req service.RenderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, err.Error())
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	rendered, err := h.promptSvc.RenderPrompt(c.Request.Context(), actor, id, &req)
	if err != nil {
		writePromptError(c, err)
		return
	}

	result.OK(c, rendered)
}

func parseID(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}

func queryInt(c *gin.Context, key string, def int) int {
	v := c.Query(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func writePromptError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrPromptNotFound) {
		result.NotFound(c, "prompt not found")
		return
	}
	if errors.Is(err, service.ErrNoPermission) {
		result.Forbidden(c, err.Error())
		return
	}
	if errors.Is(err, service.ErrMissingRequired) {
		result.BadRequest(c, err.Error())
		return
	}
	result.InternalError(c, err.Error())
}

// Ensure domain.Prompt is used in handler responses
var _ = (*domain.Prompt)(nil)
