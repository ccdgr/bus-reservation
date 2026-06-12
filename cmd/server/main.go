package main

import (
	"log"

	"github.com/ccdgr/bus-reservation/config"
	"github.com/ccdgr/bus-reservation/pkg/database"
	"github.com/ccdgr/bus-reservation/pkg/mq"
	"github.com/rabbitmq/amqp091-go"
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
	// 1. Load Config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Printf("Warning: failed to load config.yaml, using example or default values: %v", err)
		// Fallback or exit depending on strategy
		cfg, _ = config.LoadConfig("config.yaml.example")
	}

	// 2. Init MySQL
	db, err := database.NewMySQL(cfg.MySQL.DSN)
	if err != nil {
		log.Fatalf("failed to connect mysql: %v", err)
	}

	// 3. Init Redis
	rdb, err := database.NewRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	// 4. Init RabbitMQ
	conn, err := mq.NewRabbitMQ(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("failed to connect rabbitmq: %v", err)
	}
	defer conn.Close()

	container := &Container{
		Config: cfg,
		DB:     db,
		RDB:    rdb,
		MQ:     conn,
	}

	log.Printf("Application initialized successfully on port %s", container.Config.Server.Port)

	// TODO: Register repositories, usecases, and handlers
}
