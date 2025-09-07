package config

import (
	"expense-management-system/internal/auth"
	"expense-management-system/internal/db"
	"expense-management-system/internal/delivery/http"
	"expense-management-system/internal/delivery/http/middleware"
	"expense-management-system/internal/delivery/http/route"
	"expense-management-system/internal/messaging"
	"expense-management-system/internal/metrics"
	"expense-management-system/internal/repository"
	"expense-management-system/internal/usecase"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	redisrate "github.com/go-redis/redis_rate/v10"
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
	metrics.Init()

	allowedOrigins := strings.Split(cfg.Config.CorsAllowOrigins, ",")
	rateLimiter := redisrate.NewLimiter(cfg.RedisClient)

	commonMiddlewares := []gin.HandlerFunc{
		requestid.New(),
		middleware.NewRequestLoggerMiddleware(cfg.Log),
		middleware.NewRecoverMiddleware(cfg.Log),
		middleware.NewErrorMiddleware(cfg.Log),
		middleware.NewCorsMiddleware(allowedOrigins),
		middleware.NewLimiterMiddleware(rateLimiter),
	}

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
	approvalController := http.NewApprovalController(cfg.Log, cfg.Validate, approvalUsecase)

	routeCfg := route.RouteConfig{
		App:                cfg.App,
		CommonMiddlewares:  commonMiddlewares,
		AuthMiddlware:      authMiddleware,
		AuthController:     authController,
		UserController:     userController,
		ExpenseController:  expenseController,
		ApprovalController: approvalController,
	}
	routeCfg.Setup()
}
