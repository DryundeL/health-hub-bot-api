#!/bin/sh
set -e

# Entrypoint скрипт для приложения HealthHub Bot API
# Выполняет предварительные проверки перед запуском приложения

echo "=== HealthHub Bot API Entrypoint ==="

# Функция для проверки доступности PostgreSQL
wait_for_postgres() {
    if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ]; then
        echo "Warning: DB_HOST or DB_PORT not set, skipping PostgreSQL check"
        return 0
    fi

    echo "Waiting for PostgreSQL to be ready..."
    max_attempts=30
    attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if nc -z "$DB_HOST" "$DB_PORT" 2>/dev/null; then
            echo "PostgreSQL is ready!"
            return 0
        fi
        attempt=$((attempt + 1))
        echo "Attempt $attempt/$max_attempts: PostgreSQL not ready, waiting 1 second..."
        sleep 1
    done

    echo "Error: PostgreSQL is not available after $max_attempts attempts"
    return 1
}

# Функция для проверки доступности Redis
wait_for_redis() {
    if [ -z "$REDIS_HOST" ] || [ -z "$REDIS_PORT" ]; then
        echo "Warning: REDIS_HOST or REDIS_PORT not set, skipping Redis check"
        return 0
    fi

    echo "Waiting for Redis to be ready..."
    max_attempts=30
    attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if nc -z "$REDIS_HOST" "$REDIS_PORT" 2>/dev/null; then
            echo "Redis is ready!"
            return 0
        fi
        attempt=$((attempt + 1))
        echo "Attempt $attempt/$max_attempts: Redis not ready, waiting 1 second..."
        sleep 1
    done

    echo "Warning: Redis is not available after $max_attempts attempts, but continuing..."
    return 0
}

# Функция для проверки критических переменных окружения
check_required_env() {
    echo "Checking required environment variables..."
    
    required_vars="DB_HOST DB_PORT DB_USER DB_PASSWORD DB_NAME"
    missing_vars=""

    for var in $required_vars; do
        eval value=\$$var
        if [ -z "$value" ]; then
            missing_vars="$missing_vars $var"
        fi
    done

    if [ -n "$missing_vars" ]; then
        echo "Error: Missing required environment variables:$missing_vars"
        return 1
    fi

    echo "✓ All required environment variables are set"
    return 0
}

# Функция для создания необходимых директорий
ensure_directories() {
    echo "Ensuring required directories exist..."
    
    mkdir -p /app/storage
    
    echo "✓ Directories ready"
}

# Функция для логирования информации о запуске
log_startup_info() {
    echo "=== Startup Information ==="
    echo "DB_HOST: ${DB_HOST:-not set}"
    echo "DB_PORT: ${DB_PORT:-5432}"
    echo "DB_NAME: ${DB_NAME:-healthhub}"
    echo "REDIS_HOST: ${REDIS_HOST:-not set}"
    echo "REDIS_PORT: ${REDIS_PORT:-6379}"
    echo "PORT: ${PORT:-8080}"
    echo "Working directory: $(pwd)"
    echo "=========================="
}

# Основная логика запуска
main() {
    # Логируем информацию о запуске
    log_startup_info

    # Проверяем критичные переменные окружения
    if ! check_required_env; then
        echo "Error: Environment validation failed"
        exit 1
    fi

    # Ожидаем готовности PostgreSQL (только если запускается приложение, не мигратор)
    if [ "$1" = "./server" ] || [ -z "$1" ]; then
        wait_for_postgres || {
            echo "Warning: PostgreSQL check failed, but continuing..."
        }
        
        # Ожидаем готовности Redis (опционально)
        wait_for_redis
    fi

    # Создаем необходимые директории
    ensure_directories

    echo "=== Starting application ==="
    
    # Выполняем переданную команду (по умолчанию запуск приложения)
    exec "$@"
}

# Запускаем основную логику
main "$@"

