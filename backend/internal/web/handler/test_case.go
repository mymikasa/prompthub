package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mymikasa/prompthub/internal/service"
	"github.com/mymikasa/prompthub/internal/web/result"
	"github.com/mymikasa/prompthub/pkg/ctxutil"
)

type TestCaseHandler struct {
	testCaseSvc *service.TestCaseService
}

func NewTestCaseHandler(testCaseSvc *service.TestCaseService) *TestCaseHandler {
	return &TestCaseHandler{testCaseSvc: testCaseSvc}
}

func (h *TestCaseHandler) Register(rg *gin.RouterGroup) {
	prompts := rg.Group("/api/prompts")
	prompts.GET("/:id/test-cases", h.List)
	prompts.POST("/:id/test-cases", h.Create)
	prompts.PATCH("/:id/test-cases/:testCaseId", h.Update)
	prompts.DELETE("/:id/test-cases/:testCaseId", h.Delete)
}

func (h *TestCaseHandler) List(c *gin.Context) {
	promptID, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid prompt id")
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	cases, err := h.testCaseSvc.List(c.Request.Context(), actor, promptID)
	if err != nil {
		writePromptError(c, err)
		return
	}

	result.OK(c, cases)
}

func (h *TestCaseHandler) Create(c *gin.Context) {
	promptID, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid prompt id")
		return
	}

	var req service.CreateTestCaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, err.Error())
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	tc, err := h.testCaseSvc.Create(c.Request.Context(), actor, promptID, &req)
	if err != nil {
		writeTestCaseError(c, err)
		return
	}

	result.Created(c, tc)
}

func (h *TestCaseHandler) Update(c *gin.Context) {
	promptID, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid prompt id")
		return
	}

	testCaseID, err := strconv.ParseInt(c.Param("testCaseId"), 10, 64)
	if err != nil {
		result.BadRequest(c, "invalid test case id")
		return
	}

	var req service.UpdateTestCaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, err.Error())
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	tc, err := h.testCaseSvc.Update(c.Request.Context(), actor, promptID, testCaseID, &req)
	if err != nil {
		writeTestCaseError(c, err)
		return
	}

	result.OK(c, tc)
}

func (h *TestCaseHandler) Delete(c *gin.Context) {
	promptID, err := parseID(c)
	if err != nil {
		result.BadRequest(c, "invalid prompt id")
		return
	}

	testCaseID, err := strconv.ParseInt(c.Param("testCaseId"), 10, 64)
	if err != nil {
		result.BadRequest(c, "invalid test case id")
		return
	}

	actor := ctxutil.ActorFromCtx(c.Request.Context())
	if err := h.testCaseSvc.Delete(c.Request.Context(), actor, promptID, testCaseID); err != nil {
		writeTestCaseError(c, err)
		return
	}

	result.OK(c, nil)
}

func writeTestCaseError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrTestCaseNotFound) {
		result.NotFound(c, "test case not found")
		return
	}
	if errors.Is(err, service.ErrPromptNotFound) {
		result.NotFound(c, "prompt not found")
		return
	}
	if errors.Is(err, service.ErrNoPermission) {
		result.Forbidden(c, err.Error())
		return
	}
	result.InternalError(c, err.Error())
}
