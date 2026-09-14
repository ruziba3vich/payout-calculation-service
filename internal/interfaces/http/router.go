package http

import (
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Health  gin.HandlerFunc
	Auth    AuthRoutes
	Courier CourierRoutes
	Order   OrderRoutes
	Payout  PayoutRoutes
}

type Middleware struct {
	Auth          gin.HandlerFunc
	AdminOnly     gin.HandlerFunc
	SelfOrAdminID gin.HandlerFunc
}

type AuthRoutes struct {
	AdminLogin, CourierLogin gin.HandlerFunc
}

type CourierRoutes struct {
	Create, Get, List, Update, Delete gin.HandlerFunc
}

type OrderRoutes struct {
	Create, Get, List, UpdateStatus, Delete gin.HandlerFunc
}

type PayoutRoutes struct {
	Calculate, Get, List, ListByCourier gin.HandlerFunc
}

func NewRouter(h Handlers, m Middleware) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", h.Health)

	auth := r.Group("/auth")
	{
		auth.POST("/admin/login", h.Auth.AdminLogin)
		auth.POST("/courier/login", h.Auth.CourierLogin)
	}

	api := r.Group("", m.Auth)

	couriers := api.Group("/couriers")
	{
		couriers.POST("", m.AdminOnly, h.Courier.Create)
		couriers.GET("", m.AdminOnly, h.Courier.List)
		couriers.GET("/:id", m.SelfOrAdminID, h.Courier.Get)
		couriers.PATCH("/:id", m.AdminOnly, h.Courier.Update)
		couriers.DELETE("/:id", m.AdminOnly, h.Courier.Delete)
		couriers.GET("/:id/payouts", m.SelfOrAdminID, h.Payout.ListByCourier)
	}

	orders := api.Group("/orders")
	{
		orders.POST("", m.AdminOnly, h.Order.Create)
		orders.GET("", h.Order.List)
		orders.GET("/:id", h.Order.Get)
		orders.PATCH("/:id/status", h.Order.UpdateStatus)
		orders.DELETE("/:id", m.AdminOnly, h.Order.Delete)
	}

	payouts := api.Group("/payouts")
	{
		payouts.POST("/calculate", m.AdminOnly, h.Payout.Calculate)
		payouts.GET("", m.AdminOnly, h.Payout.List)
		payouts.GET("/:id", h.Payout.Get)
	}

	return r
}
