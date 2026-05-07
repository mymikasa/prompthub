package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const cookieName = "prompthub_session"

type Session struct {
	UserID      int64 `json:"uid"`
	WorkspaceID int64 `json:"wid"`
	ExpiresAt   int64 `json:"exp"`
}

func Sign(sess *Session, secret string) (string, error) {
	payload, err := json.Marshal(sess)
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(encoded))
	sig := fmt.Sprintf("%x", mac.Sum(nil))
	return encoded + "." + sig, nil
}

func Verify(token string, secret string) (*Session, error) {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid token format")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(parts[0]))
	expected := fmt.Sprintf("%x", mac.Sum(nil))
	if !hmac.Equal([]byte(parts[1]), []byte(expected)) {
		return nil, fmt.Errorf("invalid signature")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}

	var sess Session
	if err := json.Unmarshal(payload, &sess); err != nil {
		return nil, err
	}
	if time.Now().Unix() > sess.ExpiresAt {
		return nil, fmt.Errorf("session expired")
	}
	return &sess, nil
}

func SetCookie(c *gin.Context, sess *Session, secret string, isDev bool) {
	token, err := Sign(sess, secret)
	if err != nil {
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   !isDev,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func GetFromCookie(c *gin.Context, secret string) (*Session, error) {
	cookie, err := c.Cookie(cookieName)
	if err != nil {
		return nil, fmt.Errorf("cookie not found")
	}
	return Verify(cookie, secret)
}

func UserIDFromContext(c *gin.Context) int64 {
	v, _ := c.Get("user_id")
	id, _ := v.(int64)
	return id
}

func WorkspaceIDFromContext(c *gin.Context) int64 {
	v, _ := c.Get("workspace_id")
	id, _ := v.(int64)
	return id
}

func ParseUserID(str string) int64 {
	id, _ := strconv.ParseInt(str, 10, 64)
	return id
}
