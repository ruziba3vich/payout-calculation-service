package http

import (
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Health  gin.HandlerFunc
	Courier CourierRoutes
	Order   OrderRoutes
	Payout  PayoutRoutes
}

type CourierRoutes struct {
	Create, Get, List, Update, Delete gin.HandlerFunc
}

type OrderRoutes struct {
	Create, Get, List, UpdateStatus, Delete gin.HandlerFunc
}

type PayoutRoutes struct {
	Get, List, ListByCourier gin.HandlerFunc
}

func NewRouter(h Handlers) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", h.Health)

	couriers := r.Group("/couriers")
	{
		couriers.POST("", h.Courier.Create)
		couriers.GET("", h.Courier.List)
		couriers.GET("/:id", h.Courier.Get)
		couriers.PATCH("/:id", h.Courier.Update)
		couriers.DELETE("/:id", h.Courier.Delete)
		couriers.GET("/:id/payouts", h.Payout.ListByCourier)
	}

	orders := r.Group("/orders")
	{
		orders.POST("", h.Order.Create)
		orders.GET("", h.Order.List)
		orders.GET("/:id", h.Order.Get)
		orders.PATCH("/:id/status", h.Order.UpdateStatus)
		orders.DELETE("/:id", h.Order.Delete)
	}

	payouts := r.Group("/payouts")
	{
		payouts.GET("", h.Payout.List)
		payouts.GET("/:id", h.Payout.Get)
	}

	return r
}
