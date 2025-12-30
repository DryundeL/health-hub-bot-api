package config

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/gorm/logger"
)

// Config представляет конфигурацию приложения
type Config struct {
	// Database
	Database DatabaseConfig

	// Server
	Server ServerConfig

	// Telegram
	Telegram TelegramConfig

	// Storage
	Storage StorageConfig
}

// DatabaseConfig представляет конфигурацию базы данных
type DatabaseConfig struct {
	// Параметры подключения (можно задать через DATABASE_URL или отдельные переменные)
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Timezone string
	
	// Полная строка подключения (если задана DATABASE_URL)
	URL string
	
	// Настройки пула соединений
	MaxOpenConns int
	MaxIdleConns int
	LogLevel     logger.LogLevel
}

// ServerConfig представляет конфигурацию сервера
type ServerConfig struct {
	Port string
	Host string
}

// TelegramConfig представляет конфигурацию Telegram
type TelegramConfig struct {
	BotToken string
}

// StorageConfig представляет конфигурацию хранилища файлов
type StorageConfig struct {
	Type string // "local" или "s3"
	Path string // путь для локального хранилища

	// S3 конфигурация (если Type == "s3")
	S3AccessKeyID     string
	S3SecretAccessKey string
	S3Region          string
	S3Bucket          string
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	cfg := &Config{}

	// Database
	maxOpenConns := getEnvInt("DATABASE_MAX_OPEN_CONNS", 25)
	maxIdleConns := getEnvInt("DATABASE_MAX_IDLE_CONNS", 5)
	logLevelStr := os.Getenv("DATABASE_LOG_LEVEL")
	logLevel := logger.Info
	if logLevelStr == "silent" {
		logLevel = logger.Silent
	} else if logLevelStr == "error" {
		logLevel = logger.Error
	} else if logLevelStr == "warn" {
		logLevel = logger.Warn
	}

	// Поддержка DATABASE_URL для обратной совместимости
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		cfg.Database = DatabaseConfig{
			URL:          dbURL,
			MaxOpenConns: maxOpenConns,
			MaxIdleConns: maxIdleConns,
			LogLevel:     logLevel,
		}
	} else {
		// Использование отдельных переменных окружения
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")
		sslMode := getEnv("DB_SSLMODE", "disable")
		timezone := getEnv("DB_TIMEZONE", "UTC")

		if user == "" || password == "" || name == "" {
			return nil, fmt.Errorf("either DATABASE_URL or DB_USER, DB_PASSWORD, and DB_NAME environment variables must be set")
		}

		cfg.Database = DatabaseConfig{
			Host:         host,
			Port:         port,
			User:         user,
			Password:     password,
			Name:         name,
			SSLMode:      sslMode,
			Timezone:     timezone,
			MaxOpenConns: maxOpenConns,
			MaxIdleConns: maxIdleConns,
			LogLevel:     logLevel,
		}
	}

	// Server
	cfg.Server = ServerConfig{
		Port: getEnv("PORT", "8080"),
		Host: getEnv("HOST", ""),
	}

	// Telegram
	cfg.Telegram = TelegramConfig{
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}

	// Storage
	storageType := getEnv("STORAGE_TYPE", "local")
	cfg.Storage = StorageConfig{
		Type:              storageType,
		Path:              getEnv("STORAGE_PATH", "./storage"),
		S3AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		S3SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3Region:          os.Getenv("AWS_REGION"),
		S3Bucket:          os.Getenv("S3_BUCKET"),
	}

	return cfg, nil
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt возвращает значение переменной окружения как int или значение по умолчанию
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

