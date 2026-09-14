package handler

import (
	"errors"
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

func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	t, err := h.svc.LoginAdmin(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.loginError(c, err)
		return
	}

	c.JSON(http.StatusOK, toTokenResponse(t))
}

func (h *AuthHandler) CourierLogin(c *gin.Context) {
	var req courierLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	t, err := h.svc.LoginCourier(c.Request.Context(), req.Phone, req.Password)
	if err != nil {
		h.loginError(c, err)
		return
	}

	c.JSON(http.StatusOK, toTokenResponse(t))
}

func (h *AuthHandler) loginError(c *gin.Context, err error) {
	if errors.Is(err, auth.ErrInvalidCredentials) {
		httpx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	httpx.Error(c, http.StatusInternalServerError, err.Error())
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
