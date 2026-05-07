package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mymikasa/prompthub/internal/service"
	"github.com/mymikasa/prompthub/internal/web/result"
	"github.com/mymikasa/prompthub/pkg/ctxutil"
)

type PromptRunHandler struct {
	runSvc *service.PromptRunService
}

func NewPromptRunHandler(runSvc *service.PromptRunService) *PromptRunHandler {
	return &PromptRunHandler{runSvc: runSvc}
}

func (h *PromptRunHandler) Register(rg *gin.RouterGroup) {
	prompts := rg.Group("/api/prompts")
	prompts.POST("/:id/run", h.Run)
	prompts.GET("/:id/runs", h.List)

	rg.GET("/api/runs", h.ListAll)
}

func (h *PromptRunHandler) Run(c *gin.Context) {
	promptID, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid prompt id")
		return
	}

	var req service.RunPromptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, err.Error())
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	runResult, err := h.runSvc.Run(c.Request.Context(), actor, promptID, &req)
	if err != nil {
		if err == service.ErrRunFailed {
			result.InternalError(c, "prompt run failed, please check provider configuration")
			return
		}
		writePromptError(c, err)
		return
	}

	result.OK(c, runResult)
}

func (h *PromptRunHandler) List(c *gin.Context) {
	promptID, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid prompt id")
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "pageSize", 20)
	status := c.Query("status")
	model := c.Query("model")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	runs, total, err := h.runSvc.ListRuns(c.Request.Context(), actor, promptID, page, pageSize, status, model, startDate, endDate)
	if err != nil {
		writePromptError(c, err)
		return
	}

	result.OK(c, gin.H{
		"items":    runs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *PromptRunHandler) ListAll(c *gin.Context) {
	actor := ctxutil.ActorFromCtx(c.Request.Context())
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "pageSize", 20)
	status := c.Query("status")
	model := c.Query("model")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	runs, total, err := h.runSvc.ListAllRuns(c.Request.Context(), actor, page, pageSize, status, model, startDate, endDate)
	if err != nil {
		result.InternalError(c, err.Error())
		return
	}

	result.OK(c, gin.H{
		"items":    runs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}
