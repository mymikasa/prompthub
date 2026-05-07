package handler

import (
	"errors"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mymikasa/prompthub/internal/service"
	"github.com/mymikasa/prompthub/internal/web/result"
	"github.com/mymikasa/prompthub/pkg/session"
)

type AuthHandler struct {
	authService   *service.AuthService
	sessionSecret string
	isDev         bool
}

func NewAuthHandler(authService *service.AuthService, sessionSecret string, isDev bool) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		sessionSecret: sessionSecret,
		isDev:         isDev,
	}
}

func (h *AuthHandler) Register(public *gin.RouterGroup, authenticated *gin.RouterGroup) {
	auth := public.Group("/api/auth")
	auth.POST("/signup", h.Signup)
	auth.POST("/login", h.Login)
	auth.POST("/logout", h.Logout)

	authenticated.GET("/api/me", h.Me)
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req service.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, err.Error())
		return
	}

	user, err := h.authService.Signup(&req)
	if err != nil {
		slog.Error("signup failed", slog.String("error", err.Error()))
		result.BadRequest(c, err.Error())
		return
	}

	ws, err := h.authService.GetUserWorkspace(user.ID)
	if err != nil {
		result.InternalError(c, "failed to get workspace")
		return
	}

	sess := &session.Session{
		UserID:      user.ID,
		WorkspaceID: ws.ID,
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	session.SetCookie(c, sess, h.sessionSecret, h.isDev)

	result.Created(c, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"name": user.Nickname,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, err.Error())
		return
	}

	user, err := h.authService.Login(&req)
	if err != nil {
		result.Unauthorized(c, err.Error())
		return
	}

	ws, err := h.authService.GetUserWorkspace(user.ID)
	if err != nil {
		result.InternalError(c, "failed to get workspace")
		return
	}

	sess := &session.Session{
		UserID:      user.ID,
		WorkspaceID: ws.ID,
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	session.SetCookie(c, sess, h.sessionSecret, h.isDev)

	result.OK(c, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"name": user.Nickname,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session.ClearCookie(c)
	result.OK(c, nil)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := session.UserIDFromContext(c)
	if userID == 0 {
		result.Unauthorized(c, "please login")
		return
	}

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			result.NotFound(c, "user not found")
			return
		}
		result.InternalError(c, "failed to get user")
		return
	}

	workspaceID := session.WorkspaceIDFromContext(c)

	result.OK(c, gin.H{
		"id":           user.ID,
		"email":        user.Email,
		"nickname":     user.Nickname,
		"workspace_id": workspaceID,
	})
}
