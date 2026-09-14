package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	authapp "github.com/ruziba3vich/payout-calculation-service/internal/application/auth"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/auth"
	httpx "github.com/ruziba3vich/payout-calculation-service/internal/interfaces/http"
)

type AuthHandler struct {
	svc *authapp.Service
}

func NewAuthHandler(svc *authapp.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type adminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type courierLoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type tokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
	Role        auth.Role `json:"role"`
	SubjectID   uuid.UUID `json:"subject_id"`
}

// AdminLogin godoc
// @Summary  Admin login
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    body body adminLoginRequest true "credentials"
// @Success  200 {object} tokenResponse
// @Failure  400 {object} ErrorResponse
// @Failure  401 {object} ErrorResponse
// @Router   /auth/admin/login [post]
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	t, err := h.svc.LoginAdmin(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, toTokenResponse(t))
}

// CourierLogin godoc
// @Summary  Courier login
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    body body courierLoginRequest true "credentials"
// @Success  200 {object} tokenResponse
// @Failure  400 {object} ErrorResponse
// @Failure  401 {object} ErrorResponse
// @Router   /auth/courier/login [post]
func (h *AuthHandler) CourierLogin(c *gin.Context) {
	var req courierLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	t, err := h.svc.LoginCourier(c.Request.Context(), req.Phone, req.Password)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, toTokenResponse(t))
}

func toTokenResponse(t authapp.Token) tokenResponse {
	return tokenResponse{
		AccessToken: t.AccessToken,
		TokenType:   "Bearer",
		ExpiresAt:   t.ExpiresAt,
		Role:        t.Role,
		SubjectID:   t.SubjectID,
	}
}
