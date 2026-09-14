package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	courierapp "github.com/ruziba3vich/payout-calculation-service/internal/application/courier"
	httpx "github.com/ruziba3vich/payout-calculation-service/internal/interfaces/http"
)

type CourierHandler struct {
	svc *courierapp.Service
}

func NewCourierHandler(svc *courierapp.Service) *CourierHandler {
	return &CourierHandler{svc: svc}
}

type createCourierRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	HiredAt  string `json:"hired_at" binding:"required"`
}

type updateCourierRequest struct {
	FullName *string `json:"full_name"`
	Phone    *string `json:"phone"`
	Password *string `json:"password" binding:"omitempty,min=6"`
	HiredAt  *string `json:"hired_at"`
	IsActive *bool   `json:"is_active"`
}

func (h *CourierHandler) Create(c *gin.Context) {
	var req createCourierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	hiredAt, err := time.Parse("2006-01-02", req.HiredAt)
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "hired_at must be YYYY-MM-DD")
		return
	}

	courier, err := h.svc.Create(c.Request.Context(), courierapp.CreateInput{
		FullName: req.FullName,
		Phone:    req.Phone,
		Password: req.Password,
		HiredAt:  hiredAt,
	})
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toCourierResponse(courier))
}

func (h *CourierHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	courier, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, toCourierResponse(courier))
}

func (h *CourierHandler) List(c *gin.Context) {
	limit, offset := httpx.Pagination(c)

	items, total, err := h.svc.List(c.Request.Context(), courierapp.ListInput{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		writeError(c, err)
		return
	}

	httpx.List(c, toCourierResponses(items), total)
}

func (h *CourierHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateCourierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var hiredAt *time.Time
	if req.HiredAt != nil {
		t, err := time.Parse("2006-01-02", *req.HiredAt)
		if err != nil {
			httpx.Error(c, http.StatusBadRequest, "hired_at must be YYYY-MM-DD")
			return
		}
		hiredAt = &t
	}

	courier, err := h.svc.Update(c.Request.Context(), courierapp.UpdateInput{
		ID:       id,
		FullName: req.FullName,
		Phone:    req.Phone,
		Password: req.Password,
		HiredAt:  hiredAt,
		IsActive: req.IsActive,
	})
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, toCourierResponse(courier))
}

func (h *CourierHandler) Delete(c *gin.Context) {
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
