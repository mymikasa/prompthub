package middleware

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo"
	"github.com/mymikasa/prompthub/internal/web/result"
	"github.com/mymikasa/prompthub/pkg/ctxutil"
	"github.com/mymikasa/prompthub/pkg/session"
	"gorm.io/gorm"
)

func Auth(sessionSecret string, userRepo repo.UserRepo, workspaceRepo repo.WorkspaceRepo) gin.HandlerFunc {
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

		member, err := workspaceRepo.FindMemberByWorkspaceAndUser(c.Request.Context(), sess.WorkspaceID, user.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.Forbidden(c, "not a member of this workspace")
				return
			}
			result.InternalError(c, "failed to check workspace membership")
			return
		}

		actor := &domain.Actor{
			UserID:      user.ID,
			WorkspaceID: sess.WorkspaceID,
			Role:        member.Role,
		}

		c.Set("actor", actor)
		c.Set("user_id", user.ID)
		c.Set("workspace_id", sess.WorkspaceID)
		c.Set("user", user)

		ctx := ctxutil.WithActor(c.Request.Context(), actor)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
