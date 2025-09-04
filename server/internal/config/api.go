package config

import (
	"expense-management-system/internal/auth"
	"expense-management-system/internal/db"
	"expense-management-system/internal/delivery/http"
	"expense-management-system/internal/delivery/http/middleware"
	"expense-management-system/internal/delivery/http/route"
	"expense-management-system/internal/repository"
	"expense-management-system/internal/usecase"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ApiConfig struct {
	DB       *pgxpool.Pool
	TX       db.Transactioner
	App      *gin.Engine
	Log      *zap.Logger
	Validate *validator.Validate
	Config   *Env
	Producer *kafka.Producer
}

func NewApi(cfg *ApiConfig) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Config.RedisHost, cfg.Config.RedistPort),
		DB:   cfg.Config.RedistDB,
	})

	jwtToken := auth.NewJWTToken(cfg.Config.JWTSecretKey, time.Hour*24*time.Duration(cfg.Config.JWTExpirationDay))
	authMiddleware := middleware.NewAuthMiddleware(cfg.Log, redisClient, jwtToken)

	userRepository := repository.NewUserRepository(cfg.DB)

	authUsecase := usecase.NewAuthUsecase(cfg.Log, redisClient, jwtToken, userRepository)
	userUsecase := usecase.NewUserUsecase(cfg.Log, userRepository)

	authController := http.NewAuthController(cfg.Log, cfg.Validate, authUsecase)
	userController := http.NewUserController(cfg.Log, cfg.Validate, userUsecase)

	routeCfg := route.RouteConfig{
		App:            cfg.App,
		AuthMiddlware:  authMiddleware,
		AuthController: authController,
		UserController: userController,
	}
	routeCfg.Setup()
}
