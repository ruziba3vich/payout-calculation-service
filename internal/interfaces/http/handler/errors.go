package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/courier"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/order"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/payout"
	httpx "github.com/ruziba3vich/payout-calculation-service/internal/interfaces/http"
)

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, courier.ErrNotFound),
		errors.Is(err, order.ErrNotFound),
		errors.Is(err, payout.ErrNotFound):
		httpx.Error(c, http.StatusNotFound, err.Error())

	case errors.Is(err, courier.ErrPhoneTaken),
		errors.Is(err, payout.ErrAlreadyExists):
		httpx.Error(c, http.StatusConflict, err.Error())

	case errors.Is(err, payout.ErrInvalidPeriod),
		errors.Is(err, payout.ErrFuturePeriod):
		httpx.Error(c, http.StatusBadRequest, err.Error())

	default:
		httpx.Error(c, http.StatusInternalServerError, "internal error")
	}
}
