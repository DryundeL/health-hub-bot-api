package database

import (
	"fmt"
	"net/url"

	"github.com/health-hub-bot-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// buildDSN строит строку подключения DSN из отдельных параметров
func buildDSN(cfg config.DatabaseConfig) string {
	if cfg.URL != "" {
		return cfg.URL
	}

	// Экранирование специальных символов в пароле и имени пользователя для URL
	user := url.QueryEscape(cfg.User)
	password := url.QueryEscape(cfg.Password)
	dbName := url.QueryEscape(cfg.Name)

	// Формат postgres://user:password@host:port/dbname?sslmode=...&TimeZone=...
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, cfg.Host, cfg.Port, dbName, cfg.SSLMode)

	if cfg.Timezone != "" {
		dsn += fmt.Sprintf("&TimeZone=%s", url.QueryEscape(cfg.Timezone))
	}

	return dsn
}

// NewPostgres создаёт новое подключение к PostgreSQL
func NewPostgres(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := buildDSN(cfg)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Настройка пула соединений
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	return db, nil
}

// Close закрывает подключение к БД
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
