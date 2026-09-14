package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	orderapp "github.com/ruziba3vich/payout-calculation-service/internal/application/order"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/auth"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/order"
	httpx "github.com/ruziba3vich/payout-calculation-service/internal/interfaces/http"
	"github.com/ruziba3vich/payout-calculation-service/internal/interfaces/http/middleware"
)

type OrderHandler struct {
	svc *orderapp.Service
}

func NewOrderHandler(svc *orderapp.Service) *OrderHandler {
	return &OrderHandler{svc: svc}
}

type createOrderRequest struct {
	CourierID uuid.UUID       `json:"courier_id" binding:"required"`
	Amount    decimal.Decimal `json:"amount" binding:"required"`
}

type updateOrderStatusRequest struct {
	Status order.Status `json:"status" binding:"required,oneof=pending delivered cancelled returned"`
}

type listOrdersQuery struct {
	CourierID     string `form:"courier_id"`
	Status        string `form:"status" binding:"omitempty,oneof=pending delivered cancelled returned"`
	DeliveredFrom string `form:"delivered_from"`
	DeliveredTo   string `form:"delivered_to"`
	CreatedFrom   string `form:"created_from"`
	CreatedTo     string `form:"created_to"`
	AmountMin     string `form:"amount_min"`
	AmountMax     string `form:"amount_max"`
	SortBy        string `form:"sort_by" binding:"omitempty,oneof=amount delivered_at created_at"`
	SortDir       string `form:"sort_dir" binding:"omitempty,oneof=asc desc"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.svc.Create(c.Request.Context(), orderapp.CreateInput{
		CourierID: req.CourierID,
		Amount:    req.Amount,
	})
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toOrderResponse(o))
}

func (h *OrderHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	o, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}

	if !canAccessCourier(c, o.CourierID) {
		httpx.Error(c, http.StatusForbidden, "forbidden")
		return
	}

	c.JSON(http.StatusOK, toOrderResponse(o))
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if p, _ := middleware.GetPrincipal(c); p.Role == auth.RoleCourier {
		existing, err := h.svc.GetByID(c.Request.Context(), id)
		if err != nil {
			writeError(c, err)
			return
		}
		if existing.CourierID != p.ID {
			httpx.Error(c, http.StatusForbidden, "forbidden")
			return
		}
	}

	o, err := h.svc.UpdateStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, toOrderResponse(o))
}

func (h *OrderHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *OrderHandler) List(c *gin.Context) {
	var q listOrdersQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	limit, offset := httpx.Pagination(c)
	params := order.ListParams{
		SortBy:  q.SortBy,
		SortDir: q.SortDir,
		Limit:   limit,
		Offset:  offset,
	}

	if q.CourierID != "" {
		id, err := uuid.Parse(q.CourierID)
		if err != nil {
			httpx.Error(c, http.StatusBadRequest, "invalid courier_id")
			return
		}
		params.CourierID = &id
	}
	if p, _ := middleware.GetPrincipal(c); p.Role == auth.RoleCourier {
		params.CourierID = &p.ID
	}
	if q.Status != "" {
		s := order.Status(q.Status)
		params.Status = &s
	}

	var err error
	if params.DeliveredFrom, err = parseTimePtr(q.DeliveredFrom); err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid delivered_from")
		return
	}
	if params.DeliveredTo, err = parseTimePtr(q.DeliveredTo); err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid delivered_to")
		return
	}
	if params.CreatedFrom, err = parseTimePtr(q.CreatedFrom); err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid created_from")
		return
	}
	if params.CreatedTo, err = parseTimePtr(q.CreatedTo); err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid created_to")
		return
	}
	if params.AmountMin, err = parseDecimalPtr(q.AmountMin); err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid amount_min")
		return
	}
	if params.AmountMax, err = parseDecimalPtr(q.AmountMax); err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid amount_max")
		return
	}

	items, total, err := h.svc.List(c.Request.Context(), params)
	if err != nil {
		writeError(c, err)
		return
	}

	httpx.List(c, toOrderResponses(items), total)
}

func parseTimePtr(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func parseDecimalPtr(s string) (*decimal.Decimal, error) {
	if s == "" {
		return nil, nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func canAccessCourier(c *gin.Context, courierID uuid.UUID) bool {
	p, ok := middleware.GetPrincipal(c)
	if !ok {
		return false
	}
	return p.Role == auth.RoleAdmin || p.ID == courierID
}
