package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"
	httpx "github.com/ruziba3vich/payout-calculation-service/internal/interfaces/http"
)

// writeError maps the error kind to an http status. Internal errors are
// logged through gin and the client only sees a generic message.
func writeError(c *gin.Context, err error) {
	status := statusOf(errs.KindOf(err))

	if status == http.StatusInternalServerError {
		_ = c.Error(err)
		httpx.Error(c, status, "internal error")
		return
	}

	httpx.Error(c, status, errs.Public(err))
}

func statusOf(k errs.Kind) int {
	switch k {
	case errs.NotFound:
		return http.StatusNotFound
	case errs.Conflict:
		return http.StatusConflict
	case errs.Invalid:
		return http.StatusBadRequest
	case errs.Unauthorized:
		return http.StatusUnauthorized
	case errs.Forbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
