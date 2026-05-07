package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/mymikasa/prompthub/internal/repo"
	"github.com/mymikasa/prompthub/internal/web/result"
	"github.com/mymikasa/prompthub/pkg/session"
)

func Auth(sessionSecret string, userRepo *repo.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, err := session.GetFromCookie(c, sessionSecret)
		if err != nil {
			slog.Debug("auth: session invalid", slog.String("error", err.Error()))
			result.Unauthorized(c, "please login")
			return
		}

		user, err := userRepo.FindByID(c.Request.Context(), sess.UserID)
		if err != nil {
			result.Unauthorized(c, "user not found")
			return
		}

		c.Set("user_id", user.ID)
		c.Set("workspace_id", sess.WorkspaceID)
		c.Set("user", user)
		c.Next()
	}
}
