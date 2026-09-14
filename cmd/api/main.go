package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	courierapp "github.com/ruziba3vich/payout-calculation-service/internal/application/courier"
	orderapp "github.com/ruziba3vich/payout-calculation-service/internal/application/order"
	payoutapp "github.com/ruziba3vich/payout-calculation-service/internal/application/payout"
	"github.com/ruziba3vich/payout-calculation-service/internal/config"
	"github.com/ruziba3vich/payout-calculation-service/internal/infrastructure/persistence/postgres"
	httpx "github.com/ruziba3vich/payout-calculation-service/internal/interfaces/http"
	"github.com/ruziba3vich/payout-calculation-service/internal/interfaces/http/handler"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := postgres.New(ctx, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	courierRepo := postgres.NewCourierRepo(db)
	orderRepo := postgres.NewOrderRepo(db)
	payoutRepo := postgres.NewPayoutRepo(db)
	adjustmentRepo := postgres.NewPayoutAdjustmentRepo(db)

	courierSvc := courierapp.NewService(courierRepo)
	orderSvc := orderapp.NewService(orderRepo)
	payoutSvc := payoutapp.NewService(payoutRepo, adjustmentRepo)

	courierH := handler.NewCourierHandler(courierSvc)
	orderH := handler.NewOrderHandler(orderSvc)
	payoutH := handler.NewPayoutHandler(payoutSvc)

	router := httpx.NewRouter(httpx.Handlers{
		Health: handler.Health,
		Courier: httpx.CourierRoutes{
			Create: courierH.Create,
			Get:    courierH.Get,
			List:   courierH.List,
			Update: courierH.Update,
			Delete: courierH.Delete,
		},
		Order: httpx.OrderRoutes{
			Create:       orderH.Create,
			Get:          orderH.Get,
			List:         orderH.List,
			UpdateStatus: orderH.UpdateStatus,
			Delete:       orderH.Delete,
		},
		Payout: httpx.PayoutRoutes{
			Get:           payoutH.Get,
			List:          payoutH.List,
			ListByCourier: payoutH.ListByCourier,
		},
	})

	srv := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	go func() {
		log.Printf("listening on :%s", cfg.HTTP.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("shutdown error:", err)
		os.Exit(1)
	}
}
