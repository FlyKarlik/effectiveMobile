package http_response

import (
	"github.com/gin-gonic/gin"
)

type BaseResponse[T any] struct {
	Status bool  `json:"status"`
	Code   int   `json:"code"`
	Data   T     `json:"data,omitempty"`
	Error  error `json:"error,omitempty"`
}

func New[T any](c *gin.Context, code int, status bool, data T, err error) {
	obj := &BaseResponse[T]{
		Code:   code,
		Status: status,
		Data:   data,
		Error:  err,
	}
	c.JSON(code, obj)
}
