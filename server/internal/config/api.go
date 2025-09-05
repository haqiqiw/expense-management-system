package config

import (
	"expense-management-system/internal/auth"
	"expense-management-system/internal/db"
	"expense-management-system/internal/delivery/http"
	"expense-management-system/internal/delivery/http/middleware"
	"expense-management-system/internal/delivery/http/route"
	"expense-management-system/internal/messaging"
	"expense-management-system/internal/repository"
	"expense-management-system/internal/usecase"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ApiConfig struct {
	DB          *pgxpool.Pool
	TX          db.Transactioner
	App         *gin.Engine
	Log         *zap.Logger
	Validate    *validator.Validate
	Config      *Env
	Producer    *kafka.Producer
	RedisClient *redis.Client
}

func NewApi(cfg *ApiConfig) {
	jwtToken := auth.NewJWTToken(cfg.Config.JWTSecretKey, time.Hour*24*time.Duration(cfg.Config.JWTExpirationDay))
	authMiddleware := middleware.NewAuthMiddleware(cfg.Log, cfg.RedisClient, jwtToken)

	expenseApprovedProducer := messaging.NewExpenseApprovedProducer(
		cfg.Log,
		cfg.Producer,
		cfg.Config.KafkaTopicExpenseApproved,
	)

	userRepository := repository.NewUserRepository(cfg.DB)
	expenseRepository := repository.NewExpenseRepository(cfg.DB)
	approvalRepository := repository.NewApprovalRepository(cfg.DB)

	authUsecase := usecase.NewAuthUsecase(cfg.Log, cfg.RedisClient, jwtToken, userRepository)
	userUsecase := usecase.NewUserUsecase(cfg.Log, userRepository)
	expenseUsecase := usecase.NewExpenseUsecase(cfg.Log, expenseRepository, expenseApprovedProducer)
	approvalUsecase := usecase.NewApprovalUsecase(
		cfg.Log,
		cfg.TX,
		approvalRepository,
		expenseRepository,
		expenseApprovedProducer,
	)

	authController := http.NewAuthController(cfg.Log, cfg.Validate, authUsecase)
	userController := http.NewUserController(cfg.Log, cfg.Validate, userUsecase)
	expenseController := http.NewExpenseController(cfg.Log, cfg.Validate, expenseUsecase)
	approvalController := http.NewApprovalController(cfg.Log, approvalUsecase)

	routeCfg := route.RouteConfig{
		Logger:             cfg.Log,
		App:                cfg.App,
		AuthMiddlware:      authMiddleware,
		AuthController:     authController,
		UserController:     userController,
		ExpenseController:  expenseController,
		ApprovalController: approvalController,
	}
	routeCfg.Setup()
}
