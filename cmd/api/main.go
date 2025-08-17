package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/Blue-Davinci/one_acre_fund/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	version = "1.0.0" // set during build with -ldflags
)

type config struct {
	port int
	env  string
	api  struct {
		name   string
		author string
	}
	cors struct {
		trustedOrigins []string
	}
	redis struct {
		addr string
		db   int
	}
}

type application struct {
	config  config
	logger  *zap.Logger
	RedisDB *redis.Client
}

func main() {
	var cfg config

	// Logger
	logger, err := logger.InitJSONLogger()
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		os.Exit(1)
	}

	// Flags
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.api.name, "api-name", "OneAcre", "API Name")
	flag.StringVar(&cfg.api.author, "api-author", "Blue-Davinci", "API Author")
	flag.StringVar(&cfg.redis.addr, "redis-addr", "localhost:6379", "Redis address")
	flag.IntVar(&cfg.redis.db, "redis-db", 0, "Redis database")
	flag.Parse()

	// Redis init
	rdb, err := openRedis(cfg)
	if err != nil {
		logger.Fatal("Error connecting to Redis", zap.String("error", err.Error()))
	}

	app := &application{
		config:  cfg,
		logger:  logger,
		RedisDB: rdb,
	}

	// Log startup info
	logger.Info("Starting application",
		zap.String("version", version),
		zap.String("environment", cfg.env),
		zap.String("redis", cfg.redis.addr),
	)

	// Run server
	err = app.server()
	if err != nil {
		logger.Fatal("Error while starting server", zap.String("error", err.Error()))
	}
}

func openRedis(cfg config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.redis.addr,
		DB:   cfg.redis.db,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
