package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/health-hub-bot-api/graphql/generated"
	"github.com/health-hub-bot-api/internal/config"
	"github.com/health-hub-bot-api/internal/infrastructure/database"
	"github.com/health-hub-bot-api/internal/infrastructure/repository"
	"github.com/health-hub-bot-api/internal/presentation/graphql"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	// Инициализация подключения к PostgreSQL
	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	symptomRepo := repository.NewSymptomRepository(db)
	analysisRepo := repository.NewAnalysisRepository(db)
	medicationRepo := repository.NewMedicationRepository(db)
	intakeRepo := repository.NewIntakeRepository(db)
	doctorVisitRepo := repository.NewDoctorVisitRepository(db)

	// Инициализация resolver
	resolver := graphql.NewResolver(
		userRepo,
		symptomRepo,
		analysisRepo,
		medicationRepo,
		intakeRepo,
		doctorVisitRepo,
	)

	// Настройка GraphQL сервера
	// Используем type assertion для обхода проблемы с неэкспортированными типами introspection resolvers
	// Методы __Field, __InputValue, __Schema, __Type возвращают анонимные интерфейсы,
	// которые соответствуют неэкспортированным типам из generated пакета
	resolverRoot := interface{}(resolver).(generated.ResolverRoot)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolverRoot}))

	// Настройка HTTP маршрутов
	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	// Определение адреса сервера
	addr := ":" + cfg.Server.Port
	if cfg.Server.Host != "" {
		addr = cfg.Server.Host + ":" + cfg.Server.Port
	}

	// Создание HTTP сервера с таймаутами
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Канал для получения сигналов ОС
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		log.Printf("Server starting on %s", addr)
		log.Printf("GraphQL playground available at http://localhost%s/", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	<-sigChan
	log.Println("Shutting down server...")

	// Создание контекста с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Graceful shutdown сервера
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
