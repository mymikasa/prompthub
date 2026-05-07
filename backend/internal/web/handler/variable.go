package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mymikasa/prompthub/internal/service"
	"github.com/mymikasa/prompthub/internal/web/result"
	"github.com/mymikasa/prompthub/pkg/ctxutil"
)

type VariableHandler struct {
	variableSvc *service.VariableService
	promptSvc   *service.PromptService
}

func NewVariableHandler(variableSvc *service.VariableService, promptSvc *service.PromptService) *VariableHandler {
	return &VariableHandler{variableSvc: variableSvc, promptSvc: promptSvc}
}

func (h *VariableHandler) Register(rg *gin.RouterGroup) {
	prompts := rg.Group("/api/prompts")
	prompts.GET("/:id/variables", h.List)
	prompts.PATCH("/:id/variables/:variableId", h.Update)
}

func (h *VariableHandler) List(c *gin.Context) {
	promptID, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid prompt id")
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	if _, err := h.promptSvc.GetPrompt(c.Request.Context(), actor, promptID); err != nil {
		writePromptError(c, err)
		return
	}

	variables, err := h.variableSvc.ListVariables(c.Request.Context(), promptID)
	if err != nil {
		result.InternalError(c, err.Error())
		return
	}

	result.OK(c, variables)
}

func (h *VariableHandler) Update(c *gin.Context) {
	promptID, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid prompt id")
		return
	}

	variableID, err := strconv.ParseInt(c.Param("variableId"), 10, 64)
	if err != nil {
		result.BadRequest(c, "invalid variable id")
		return
	}

	var req service.UpdateVariableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, err.Error())
		return
	}

	v, err := h.variableSvc.UpdateVariable(c.Request.Context(), promptID, variableID, &req)
	if err != nil {
		if errors.Is(err, service.ErrVariableNotFound) {
			result.NotFound(c, "variable not found")
			return
		}
		result.InternalError(c, err.Error())
		return
	}

	result.OK(c, v)
}
