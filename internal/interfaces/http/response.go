package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type errorBody struct {
	Error string `json:"error"`
}

type listBody struct {
	Items any   `json:"items"`
	Total int64 `json:"total"`
}

func Error(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, errorBody{Error: msg})
}

func List(c *gin.Context, items any, total int64) {
	c.JSON(http.StatusOK, listBody{Items: items, Total: total})
}

func Pagination(c *gin.Context) (limit, offset int32) {
	limit, offset = 20, 0
	if v, err := strconv.Atoi(c.Query("limit")); err == nil && v > 0 && v <= 100 {
		limit = int32(v)
	}
	if v, err := strconv.Atoi(c.Query("offset")); err == nil && v >= 0 {
		offset = int32(v)
	}
	return limit, offset
}
