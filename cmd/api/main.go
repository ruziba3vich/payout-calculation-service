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

	adminapp "github.com/ruziba3vich/payout-calculation-service/internal/application/admin"
	authapp "github.com/ruziba3vich/payout-calculation-service/internal/application/auth"
	courierapp "github.com/ruziba3vich/payout-calculation-service/internal/application/courier"
	orderapp "github.com/ruziba3vich/payout-calculation-service/internal/application/order"
	payoutapp "github.com/ruziba3vich/payout-calculation-service/internal/application/payout"
	"github.com/ruziba3vich/payout-calculation-service/internal/config"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/auth"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/payout"
	jwtauth "github.com/ruziba3vich/payout-calculation-service/internal/infrastructure/auth"
	rediscache "github.com/ruziba3vich/payout-calculation-service/internal/infrastructure/cache"
	"github.com/ruziba3vich/payout-calculation-service/internal/infrastructure/persistence/postgres"
	"github.com/ruziba3vich/payout-calculation-service/internal/infrastructure/scheduler"
	httpx "github.com/ruziba3vich/payout-calculation-service/internal/interfaces/http"
	"github.com/ruziba3vich/payout-calculation-service/internal/interfaces/http/handler"
	"github.com/ruziba3vich/payout-calculation-service/internal/interfaces/http/middleware"
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

	redis, err := rediscache.NewRedis(ctx, cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}
	defer redis.Close()

	adminRepo := postgres.NewAdminRepo(db)
	courierRepo := postgres.NewCourierRepo(db)
	orderRepo := postgres.NewOrderRepo(db)
	payoutRepo := postgres.NewPayoutRepo(db)
	adjustmentRepo := postgres.NewPayoutAdjustmentRepo(db)
	transactor := postgres.NewTransactor(db)
	locker := postgres.NewLocker(db)

	tokens := jwtauth.NewManager(cfg.JWT)

	adminSvc := adminapp.NewService(adminRepo)
	authSvc := authapp.NewService(adminRepo, courierRepo, tokens)
	courierSvc := courierapp.NewService(courierRepo)
	calculator := payout.NewCalculator(payout.DefaultTiers)
	payoutCache := payoutapp.NewCache(redis, cfg.Redis.PayoutTTL, cfg.Redis.PayoutListTTL)
	payoutSvc := payoutapp.NewService(payoutRepo, adjustmentRepo, courierRepo, transactor, calculator, payoutCache)
	monthlyJob := payoutapp.NewMonthlyJob(payoutSvc, courierRepo, locker)
	orderSvc := orderapp.NewService(orderRepo, courierRepo, transactor, payoutSvc)

	if cfg.Admin.Username != "" && cfg.Admin.Password != "" {
		if err := adminSvc.EnsureExists(ctx, cfg.Admin.Username, cfg.Admin.Password); err != nil {
			log.Fatal("seed admin: ", err)
		}
	}

	authH := handler.NewAuthHandler(authSvc)
	courierH := handler.NewCourierHandler(courierSvc)
	orderH := handler.NewOrderHandler(orderSvc)
	payoutH := handler.NewPayoutHandler(payoutSvc, monthlyJob)

	router := httpx.NewRouter(httpx.Handlers{
		Health: handler.Health,
		Auth: httpx.AuthRoutes{
			AdminLogin:   authH.AdminLogin,
			CourierLogin: authH.CourierLogin,
		},
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
			Calculate:     payoutH.Calculate,
			RunJob:        payoutH.RunJob,
			Get:           payoutH.Get,
			List:          payoutH.List,
			ListByCourier: payoutH.ListByCourier,
		},
	}, httpx.Middleware{
		Auth:          middleware.Auth(tokens),
		AdminOnly:     middleware.RequireRole(auth.RoleAdmin),
		SelfOrAdminID: middleware.RequireSelfOrAdmin("id"),
	})

	sched := scheduler.New()
	if cfg.Job.Enabled {
		err := sched.Add(cfg.Job.Cron, "monthly payout", func(ctx context.Context) error {
			_, err := monthlyJob.RunPrevious(ctx)
			return err
		})
		if err != nil {
			log.Fatal(err)
		}
		sched.Start()
		log.Printf("monthly payout job scheduled: %q", cfg.Job.Cron)
	}

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

	sched.Stop(shutdownCtx)

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("shutdown error:", err)
		os.Exit(1)
	}
}
