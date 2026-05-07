package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mymikasa/prompthub/internal/service"
	"github.com/mymikasa/prompthub/internal/web/result"
	"github.com/mymikasa/prompthub/pkg/ctxutil"
)

type TagHandler struct {
	tagSvc *service.TagService
}

func NewTagHandler(tagSvc *service.TagService) *TagHandler {
	return &TagHandler{tagSvc: tagSvc}
}

func (h *TagHandler) Register(rg *gin.RouterGroup) {
	tags := rg.Group("/api/tags")
	tags.GET("", h.List)
}

func (h *TagHandler) List(c *gin.Context) {
	actor := ctxutil.ActorFromCtx(c.Request.Context())
	names, err := h.tagSvc.ListTags(c.Request.Context(), actor)
	if err != nil {
		result.InternalError(c, err.Error())
		return
	}
	result.OK(c, names)
}
