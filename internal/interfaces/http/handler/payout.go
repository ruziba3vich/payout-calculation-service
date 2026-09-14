package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	payoutapp "github.com/ruziba3vich/payout-calculation-service/internal/application/payout"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/payout"
	httpx "github.com/ruziba3vich/payout-calculation-service/internal/interfaces/http"
)

type PayoutHandler struct {
	svc *payoutapp.Service
	job *payoutapp.MonthlyJob
}

func NewPayoutHandler(svc *payoutapp.Service, job *payoutapp.MonthlyJob) *PayoutHandler {
	return &PayoutHandler{svc: svc, job: job}
}

type runJobRequest struct {
	Period string `json:"period"` // optional YYYY-MM, default previous month
}

// RunJob triggers the monthly job by hand. Same lock as the cron run.
// RunJob godoc
// @Summary  Run monthly payout job
// @Description Calculates payouts for all active couriers for one month. Default is the previous month. Same advisory lock as the cron run, a second parallel call returns skipped=true.
// @Tags     payouts
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    body body runJobRequest false "optional period"
// @Success  200 {object} payoutapp.JobResult
// @Failure  400 {object} ErrorResponse
// @Failure  401 {object} ErrorResponse
// @Failure  403 {object} ErrorResponse
// @Router   /payouts/jobs/monthly [post]
func (h *PayoutHandler) RunJob(c *gin.Context) {
	var req runJobRequest
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			httpx.Error(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	var (
		res payoutapp.JobResult
		err error
	)
	if req.Period == "" {
		res, err = h.job.RunPrevious(c.Request.Context())
	} else {
		period, perr := payout.ParsePeriod(req.Period)
		if perr != nil {
			writeError(c, perr)
			return
		}
		res, err = h.job.Run(c.Request.Context(), period)
	}
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

type calculatePayoutRequest struct {
	CourierID uuid.UUID `json:"courier_id" binding:"required"`
	Period    string    `json:"period" binding:"required"`
}

// Calculate creates the payout for courier + month. A second call for the same
// pair returns 409 with the existing payout in the body.
// Calculate godoc
// @Summary  Calculate payout for courier + month
// @Description Idempotent. A second call for the same courier and month returns 409 with the existing payout.
// @Tags     payouts
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    body body calculatePayoutRequest true "courier id and period YYYY-MM"
// @Success  201 {object} PayoutResponse
// @Failure  400 {object} ErrorResponse "bad or future period"
// @Failure  401 {object} ErrorResponse
// @Failure  403 {object} ErrorResponse
// @Failure  404 {object} ErrorResponse "courier not found"
// @Failure  409 {object} PayoutConflictResponse "already calculated"
// @Router   /payouts/calculate [post]
func (h *PayoutHandler) Calculate(c *gin.Context) {
	var req calculatePayoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	period, err := payout.ParsePeriod(req.Period)
	if err != nil {
		writeError(c, err)
		return
	}

	p, err := h.svc.Calculate(c.Request.Context(), req.CourierID, period)
	if errors.Is(err, payout.ErrAlreadyExists) {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error":  err.Error(),
			"payout": toPayoutResponse(p),
		})
		return
	}
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toPayoutResponse(p))
}

type listPayoutsQuery struct {
	CourierID  string `form:"courier_id"`
	Status     string `form:"status" binding:"omitempty,oneof=calculated paid cancelled"`
	PeriodFrom string `form:"period_from"`
	PeriodTo   string `form:"period_to"`
	SortBy     string `form:"sort_by" binding:"omitempty,oneof=period net_amount created_at"`
	SortDir    string `form:"sort_dir" binding:"omitempty,oneof=asc desc"`
}

// Get godoc
// @Summary  Get payout with adjustments
// @Tags     payouts
// @Produce  json
// @Security BearerAuth
// @Param    id path string true "payout id"
// @Success  200 {object} PayoutWithAdjustmentsResponse
// @Failure  401 {object} ErrorResponse
// @Failure  403 {object} ErrorResponse
// @Failure  404 {object} ErrorResponse
// @Router   /payouts/{id} [get]
func (h *PayoutHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	res, err := h.svc.GetWithAdjustments(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}

	if !canAccessCourier(c, res.Payout.CourierID) {
		httpx.Error(c, http.StatusForbidden, "forbidden")
		return
	}

	c.JSON(http.StatusOK, PayoutWithAdjustmentsResponse{
		PayoutResponse: toPayoutResponse(res.Payout),
		Adjustments:    toAdjustmentResponses(res.Adjustments),
		TotalAmount:    payoutapp.Total(res.Payout, res.Adjustments),
	})
}

// List godoc
// @Summary  List all payouts
// @Tags     payouts
// @Produce  json
// @Security BearerAuth
// @Param    courier_id  query string false "courier id"
// @Param    status      query string false "calculated | paid | cancelled"
// @Param    period_from query string false "YYYY-MM"
// @Param    period_to   query string false "YYYY-MM"
// @Param    sort_by     query string false "period | net_amount | created_at"
// @Param    sort_dir    query string false "asc | desc"
// @Param    limit       query int    false "page size, max 100" default(20)
// @Param    offset      query int    false "offset" default(0)
// @Success  200 {object} PayoutListResponse
// @Failure  400 {object} ErrorResponse
// @Failure  401 {object} ErrorResponse
// @Failure  403 {object} ErrorResponse
// @Router   /payouts [get]
func (h *PayoutHandler) List(c *gin.Context) {
	params, ok := h.listParams(c)
	if !ok {
		return
	}

	items, total, err := h.svc.List(c.Request.Context(), params)
	if err != nil {
		writeError(c, err)
		return
	}

	httpx.List(c, toPayoutResponses(items), total)
}

// ListByCourier godoc
// @Summary  Courier payout history
// @Tags     couriers
// @Produce  json
// @Security BearerAuth
// @Param    id          path  string true  "courier id"
// @Param    status      query string false "calculated | paid | cancelled"
// @Param    period_from query string false "YYYY-MM"
// @Param    period_to   query string false "YYYY-MM"
// @Param    sort_by     query string false "period | net_amount | created_at"
// @Param    sort_dir    query string false "asc | desc"
// @Param    limit       query int    false "page size, max 100" default(20)
// @Param    offset      query int    false "offset" default(0)
// @Success  200 {object} PayoutListResponse
// @Failure  400 {object} ErrorResponse
// @Failure  401 {object} ErrorResponse
// @Failure  403 {object} ErrorResponse
// @Router   /couriers/{id}/payouts [get]
func (h *PayoutHandler) ListByCourier(c *gin.Context) {
	courierID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	params, ok := h.listParams(c)
	if !ok {
		return
	}
	params.CourierID = &courierID

	items, total, err := h.svc.List(c.Request.Context(), params)
	if err != nil {
		writeError(c, err)
		return
	}

	httpx.List(c, toPayoutResponses(items), total)
}

func (h *PayoutHandler) listParams(c *gin.Context) (payout.ListParams, bool) {
	var q listPayoutsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return payout.ListParams{}, false
	}

	limit, offset := httpx.Pagination(c)
	params := payout.ListParams{
		SortBy:  q.SortBy,
		SortDir: q.SortDir,
		Limit:   limit,
		Offset:  offset,
	}

	if q.CourierID != "" {
		id, err := uuid.Parse(q.CourierID)
		if err != nil {
			httpx.Error(c, http.StatusBadRequest, "invalid courier_id")
			return params, false
		}
		params.CourierID = &id
	}
	if q.Status != "" {
		s := payout.Status(q.Status)
		params.Status = &s
	}
	if q.PeriodFrom != "" {
		t, err := time.Parse("2006-01", q.PeriodFrom)
		if err != nil {
			httpx.Error(c, http.StatusBadRequest, "period_from must be YYYY-MM")
			return params, false
		}
		params.PeriodFrom = &t
	}
	if q.PeriodTo != "" {
		t, err := time.Parse("2006-01", q.PeriodTo)
		if err != nil {
			httpx.Error(c, http.StatusBadRequest, "period_to must be YYYY-MM")
			return params, false
		}
		params.PeriodTo = &t
	}

	return params, true
}
