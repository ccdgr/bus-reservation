package main

import (
	"log/slog"
	"os"

	"github.com/ccdgr/bus-reservation/config"
	"github.com/ccdgr/bus-reservation/pkg/database"
	"github.com/ccdgr/bus-reservation/pkg/mq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
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
		// Fallback or exit depending on strategy
		cfg, _ = config.LoadConfig("config.yaml.example")
	}

	// 2. Init MySQL
	db, err := database.NewMySQL(cfg.MySQL.DSN)
	if err != nil {
		slog.Error("failed to connect mysql", "error", err)
		os.Exit(1)
	}

	// 3. Init Redis
	rdb, err := database.NewRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		slog.Error("failed to connect redis", "error", err)
		os.Exit(1)
	}

	// 4. Init RabbitMQ
	conn, err := mq.NewRabbitMQ(cfg.RabbitMQ.URL)
	if err != nil {
		slog.Error("failed to connect rabbitmq", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	container := &Container{
		Config: cfg,
		DB:     db,
		RDB:    rdb,
		MQ:     conn,
	}

	slog.Info("Application initialized successfully", "port", container.Config.Server.Port)

	// TODO: Register repositories, usecases, and handlers
}
