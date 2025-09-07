package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	AppName         string
	AppPort         int
	AppReadTimeout  int
	AppWriteTimeout int
	AppIdleTimeout  int

	CorsAllowOrigins string

	DBHost            string
	DBPort            string
	DBName            string
	DBUsername        string
	DBPassword        string
	DBMaxConnLifetime int
	DBMaxConn         int
	DBMinConn         int
	DBSSLMode         string
	DBConnectTimeout  int

	RedisHost  string
	RedistPort string
	RedistDB   int

	JWTSecretKey     string
	JWTExpirationDay int

	KafkaBrokerHost           string
	KafkaConsumerGroup        string
	KafkaAutoOffsetReset      string
	KafkaTopicExpenseApproved string
	KafkaMaxRetries           int
	KafkaBackoffDuration      int
	KafkaMaxExecuteDuration   int

	PaymentPartnerHost    string
	PaymentPartnerTimeout int
	PaymentLockDuration   int
}

func NewEnv() (*Env, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("failed to load env = %w", err)
		}
	}

	cfg := &Env{
		AppName:         getEnvString("APP_NAME", "expense-management"),
		AppPort:         getEnvInt("APP_PORT", 8500),
		AppReadTimeout:  getEnvInt("APP_READ_TIMEOUT", 60),
		AppWriteTimeout: getEnvInt("APP_WRITE_TIMEOUT", 60),
		AppIdleTimeout:  getEnvInt("APP_IDLE_TIMEOUT", 120),

		CorsAllowOrigins: getEnvString("CORS_ALLOW_ORIGINS", "http://localhost:5173"),

		DBHost:            getEnvString("DATABASE_HOST", "127.0.0.1"),
		DBPort:            getEnvString("DATABASE_PORT", "5432"),
		DBName:            getEnvString("DATABASE_NAME", ""),
		DBUsername:        getEnvString("DATABASE_USERNAME", ""),
		DBPassword:        getEnvString("DATABASE_PASSWORD", ""),
		DBMaxConnLifetime: getEnvInt("DATABASE_MAX_CONN_LIFETIME", 180),
		DBMaxConn:         getEnvInt("DATABASE_MAX_CONN", 10),
		DBMinConn:         getEnvInt("DATABASE_MIN_CONN", 10),
		DBSSLMode:         getEnvString("DATABASE_SSL_MODE", "disable"),
		DBConnectTimeout:  getEnvInt("DATABASE_CONNECT_TIMEOUT", 5),

		RedisHost:  getEnvString("REDIS_HOST", "127.0.0.1"),
		RedistPort: getEnvString("REDIS_PORT", "6379"),
		RedistDB:   getEnvInt("REDIS_DB", 0),

		JWTSecretKey:     getEnvString("JWT_SECRET_KEY", ""),
		JWTExpirationDay: getEnvInt("JWT_EXPIRATION_DAY", 1),

		KafkaBrokerHost:           getEnvString("KAFKA_BROKER_HOST", "127.0.0.1:9092"),
		KafkaConsumerGroup:        getEnvString("KAFKA_CONSUMER_GROUP", "expense-management"),
		KafkaAutoOffsetReset:      getEnvString("KAFKA_AUTO_OFFSET_RESET", "latest"),
		KafkaTopicExpenseApproved: getEnvString("KAFKA_TOPIC_EXPENSE_APPROVED", "expense-approved"),

		KafkaMaxRetries:         getEnvInt("KAFKA_MAX_RETRIES", 3),
		KafkaBackoffDuration:    getEnvInt("KAFKA_BACKOFF_DURATION", 1),
		KafkaMaxExecuteDuration: getEnvInt("KAFKA_MAX_EXECUTE_DURATION", 10),

		PaymentPartnerHost:    getEnvString("PAYMENT_PARTNER_HOST", "http://127.0.0.1:9500"),
		PaymentPartnerTimeout: getEnvInt("PAYMENT_PARTNER_TIMEOUT", 3),
		PaymentLockDuration:   getEnvInt("PAYMENT_LOCK_DURATION", 30),
	}

	return cfg, nil
}

func getEnvString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		pVal, err := strconv.Atoi(val)
		if err != nil {
			return defaultVal
		}
		return pVal
	}

	return defaultVal
}
