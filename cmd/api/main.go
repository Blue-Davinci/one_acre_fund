package main

import (
	"flag"
	"fmt"

	"github.com/Blue-Davinci/one_acre_fund/internal/logger"
	"go.uber.org/zap"
)

var (
	version = "1.0.0" // Version of the application, set at build time
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
}

type application struct {
	config config
	logger *zap.Logger
}

func main() {
	var cfg config
	logger, err := logger.InitJSONLogger()
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}
	// Port & env
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	// api configuration
	flag.StringVar(&cfg.api.name, "api-name", "OneAcre", "API Name")
	flag.StringVar(&cfg.api.author, "api-author", "Blue-Davinci", "API Author")
	// Parse the flags
	flag.Parse()
	// initialize our app
	app := &application{
		config: cfg,
		logger: logger,
	}
	// Initialize the server
	logger.Info("Loaded Cors Origins", zap.Strings("origins", cfg.cors.trustedOrigins))
	logger.Info("Starting application", zap.String("version", version), zap.String("environment", cfg.env))
	err = app.server()
	if err != nil {
		logger.Fatal("Error while starting server.", zap.String("error", err.Error()))
	}
}
