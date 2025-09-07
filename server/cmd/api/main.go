package main

import (
	"context"
	"errors"
	"expense-management-system/internal/config"
	"expense-management-system/internal/db"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	logger, err := config.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	env, err := config.NewEnv()
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to initialize env: %+v", err))
	}

	database, err := config.NewDatabase(ctx, env)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to initialize database: %+v", err))
	}
	tx := db.NewTransactioner(database)

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", env.RedisHost, env.RedistPort),
		DB:   env.RedistDB,
	})

	producer, err := config.NewKafkaProducer(env, logger)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to initialize producer: %+v", err))
	}

	validate := config.NewValidator()
	app := config.NewGin(logger)

	config.NewApi(&config.ApiConfig{
		DB:          database,
		TX:          tx,
		App:         app,
		Log:         logger,
		Validate:    validate,
		Config:      env,
		Producer:    producer,
		RedisClient: redisClient,
	})

	serverAddr := fmt.Sprintf(":%d", env.AppPort)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      app,
		ReadTimeout:  time.Duration(env.AppReadTimeout) * time.Second,
		WriteTimeout: time.Duration(env.AppWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(env.AppIdleTimeout) * time.Second,
	}
	serverErrCh := make(chan error, 1)

	go func() {
		logger.Info(fmt.Sprintf("starting server at port %d", env.AppPort))
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrCh:
		logger.Info(fmt.Sprintf("server failed to listen and serve: %+v", err))
	case sig := <-sigCh:
		logger.Info("server received stop signal", zap.String("signal", sig.String()))
	}

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		logger.Info(fmt.Sprintf("server failed to shutdown: %+v", err))
		os.Exit(1)
	}
	logger.Info("server exited properly")
}
