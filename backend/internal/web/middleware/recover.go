package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered",
					slog.Any("error", err),
					slog.String("stack", string(debug.Stack())),
					slog.String("path", c.Request.URL.Path),
					slog.String("method", c.Request.Method),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    50000,
					"message": "internal server error",
				})
			}
		}()
		c.Next()
	}
}
