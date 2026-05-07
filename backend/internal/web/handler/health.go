package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mymikasa/prompthub/internal/web/result"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Register(r *gin.RouterGroup) {
	r.GET("/healthz", h.Check)
}

func (h *HealthHandler) Check(c *gin.Context) {
	db, err := h.db.DB()
	if err != nil {
		result.InternalError(c, "database unavailable")
		return
	}
	if err := db.Ping(); err != nil {
		result.InternalError(c, "database unavailable")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data": gin.H{
			"status": "healthy",
		},
	})
}
