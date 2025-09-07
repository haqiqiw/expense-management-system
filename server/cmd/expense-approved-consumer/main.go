package main

import (
	"context"
	"expense-management-system/internal/config"
	"expense-management-system/internal/db"
	"expense-management-system/internal/delivery/messaging"
	"expense-management-system/internal/httpclient"
	"expense-management-system/internal/repository"
	"expense-management-system/internal/usecase"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
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

	kafkaConsumer, err := config.NewKafkaConsumer(env, logger)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to initialize consumer: %+v", err))
	}

	paymentPartnerClient := httpclient.NewClient(
		env.PaymentPartnerHost,
		time.Second*time.Duration(env.PaymentPartnerTimeout),
	)

	expenseRepository := repository.NewExpenseRepository(database)
	paymentPartnerRepository := repository.NewPaymentPartnerRepository(paymentPartnerClient)
	paymentProcessorUsecase := usecase.NewPaymentProcessorUsecase(
		logger,
		redisClient,
		tx,
		expenseRepository,
		paymentPartnerRepository,
		env.PaymentLockDuration,
	)

	expenseHandler := messaging.NewExpenseApprovedHandler(logger, paymentProcessorUsecase)

	consumerCfg := &messaging.ConsumerConfig{
		Topic:              env.KafkaTopicExpenseApproved,
		MaxRetries:         env.KafkaMaxRetries,
		BackoffDuration:    time.Second * time.Duration(env.KafkaBackoffDuration),
		MaxExecuteDuration: time.Second * time.Duration(env.KafkaMaxExecuteDuration),
	}
	expenseConsumer, err := messaging.NewConsumer(logger, kafkaConsumer, consumerCfg, expenseHandler.Consume)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to start consumer: %+v", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := expenseConsumer.Consume(ctx)
		if err != nil {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-quit:
		logger.Info("stop signal received, shutting down...", zap.String("signal", s.String()))
	case e := <-errCh:
		logger.Error("consumer error, shutting down...", zap.Error(e))
	}

	cancel()
	wg.Wait()

	logger.Info("consumer exited properly")
}
