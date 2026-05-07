package result

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "created",
		Data:    data,
	})
}

type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Abort(c *gin.Context, httpStatus int, code int, msg string) {
	c.AbortWithStatusJSON(httpStatus, Response{
		Code:    code,
		Message: msg,
	})
}

func BadRequest(c *gin.Context, msg string) {
	Abort(c, http.StatusBadRequest, 40000, msg)
}

func Unauthorized(c *gin.Context, msg string) {
	Abort(c, http.StatusUnauthorized, 40100, msg)
}

func Forbidden(c *gin.Context, msg string) {
	Abort(c, http.StatusForbidden, 40300, msg)
}

func NotFound(c *gin.Context, msg string) {
	Abort(c, http.StatusNotFound, 40400, msg)
}

func InternalError(c *gin.Context, msg string) {
	Abort(c, http.StatusInternalServerError, 50000, msg)
}
