package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ccdgr/bus-reservation/config"
	deliveryHTTP "github.com/ccdgr/bus-reservation/internal/delivery/http"
	deliveryMQ "github.com/ccdgr/bus-reservation/internal/delivery/mq"
	"github.com/ccdgr/bus-reservation/internal/repository"
	"github.com/ccdgr/bus-reservation/internal/usecase"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/ccdgr/bus-reservation/pkg/database"
	"github.com/ccdgr/bus-reservation/pkg/mq"
)

type Container struct {
	Config *config.Config
	DB     *gorm.DB
	RDB    *redis.Client
	MQ     *amqp.Connection
}

func main() {
	// Initialize slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 1. Load Config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		slog.Warn("failed to load config.yaml, using example or default values", "error", err)
		cfg, _ = config.LoadConfig("config.yaml.example")
	}
	if cfg.Server.JWTSecret == "" {
		cfg.Server.JWTSecret = "bus-reservation-secret-key"
	}

	// 2. Init Middleware Connections
	db, err := database.NewMySQL(cfg.MySQL.DSN)
	if err != nil {
		slog.Error("failed to connect mysql", "error", err)
		os.Exit(1)
	}

	rdb, err := database.NewRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		slog.Error("failed to connect redis", "error", err)
		os.Exit(1)
	}

	mqConn, err := mq.NewRabbitMQ(cfg.RabbitMQ.URL)
	if err != nil {
		slog.Error("failed to connect rabbitmq", "error", err)
		os.Exit(1)
	}
	defer mqConn.Close()

	ch, err := mqConn.Channel()
	if err != nil {
		slog.Error("failed to open mq channel", "error", err)
		os.Exit(1)
	}
	defer ch.Close()

	// 3. Initialize Repositories
	userRepo := repository.NewMySQLUserRepository(db)
	busRepo := repository.NewMySQLBusRepository(db)
	orderRepo := repository.NewMySQLOrderRepository(db)
	redisRepo := repository.NewRedisBusRepository(rdb)

	// 4. Initialize Usecases
	userUsecase := usecase.NewUserUsecase(userRepo, cfg.Server.JWTSecret)
	busUsecase := usecase.NewBusUsecase(busRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, busRepo, redisRepo, ch)

	// 4.5 Warmup Redis Cache (Populate stock from MySQL)
	go func() {
		slog.Info("warming up redis cache...")
		buses, err := busRepo.List(context.Background(), "", "", "")
		if err != nil {
			slog.Error("failed to list buses for warmup", "error", err)
			return
		}
		for _, b := range buses {
			if err := redisRepo.SetStock(context.Background(), b.ID, b.LeftSeat); err != nil {
				slog.Error("failed to set redis stock for bus", "bus_id", b.ID, "error", err)
			}
		}
		slog.Info("redis cache warmup completed", "count", len(buses))
	}()

	// 5. Initialize Delivery (HTTP & MQ)
	r := gin.Default()
	deliveryHTTP.NewRouter(r, userUsecase, busUsecase, orderUsecase, cfg.Server.JWTSecret)

	// Start MQ Consumer
	consumer := deliveryMQ.NewOrderConsumer(mqConn, orderRepo, busRepo)
	ctx, cancel := context.WithCancel(context.Background())
	go consumer.Start(ctx)

	// 6. Start Server (with Graceful Shutdown)
	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen: %s\n", "error", err)
		}
	}()

	slog.Info("Application initialized successfully", "port", cfg.Server.Port)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	cancel() // Stop MQ consumer
	
	ctx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown:", "error", err)
	}

	slog.Info("Server exiting")
}
