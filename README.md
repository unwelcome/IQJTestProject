# IQJ Test Task

Монолитный backend сервер на Go + fiber для управления котиками с авторизацией и загрузкой фотографий.

## Технологии

- **Backend**: Go, Fiber, Zerolog
- **Базы данных**: PostgreSQL, Redis, MinIO
- **Документация**: Swaggo/Swag
- **Контейнеризация**: Docker, Docker Compose
- **Аутентификация**: JWT tokens

## Функциональность

- JWT аутентификация и авторизация
- CRUD операции с котиками
- Загрузка и управление фотографиями котиков
- Проверка прав владения котиками
- Автоматическая документация API

## Установка и запуск

### 1. Клонирование репозитория

```bash
git clone https://github.com/unwelcome/IQJTestProject
cd IQJTestProject
```

### 2. Настройка окружения

Создайте файл `.env` в корне проекта:

```env
# Backend настройки
BACKEND_INTERNAL_HOST=0.0.0.0
BACKEND_INTERNAL_PORT=8080
BACKEND_PUBLIC_HOST=localhost
BACKEND_PUBLIC_PORT=8080

# JWT секрет
JWT_SECRET=kjmdfskjaoiwaj9fjwop3q34wstgr

# PostgreSQL
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_DB=app_db
POSTGRES_USER=postgres
POSTGRES_PASSWORD=12345678

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=redis12345678
REDIS_USER=default
REDIS_DB=0

# MinIO
MINIO_HOST=minio
MINIO_PORT=9000
MINIO_USER=minio
MINIO_PASSWORD=minio12345678
MINIO_SSL=false
```

### 3. Запуск приложения

```bash
# Запуск всех сервисов
docker compose up --build -d

# Просмотр запущенных контейнеров
docker ps

# Просмотр логов
docker logs backend
```

### 4. Проверка работоспособности

После запуска приложение будет доступно по адресу:
- **API**: http://localhost:8080/api
- **Документация Swagger**: http://localhost:8080/api/swagger/index.html
- **MinIO Console**: http://localhost:9001 (логин: minio, пароль: minio12345678)

## API Документация

После запуска проекта документация API доступна через Swagger UI:
```
http://localhost:8080/api/swagger/index.html
```

Для обновления документации (при изменении кода):
```bash
cd ./backend
swag init -o ./api/docs --dir ./cmd/api,./internal/entities,./internal/handlers
```

## Разработка

### Локальный запуск без Docker

```bash
# Убедитесь, что базы данных запущены в Docker
docker compose up postgres redis minio -d

# Запуск бекенда локально
go run cmd/api/main.go
```

### Структура проекта

```
IQJTestProject/
├── backend/
│   ├── api/docs/                   # Swagger документация
│   ├── cmd/api/                    # Точка входа приложения
│   ├── database/
│       ├── minio/                  # MinIO подключение
│       ├── postgresql/             # PostgreSQL подключение
│       ├── redisdb/                # Redis подключение
│       └── connect.go              # Подключение всех БД
│   ├── internal/
│       ├── config/                 # Конфигурация
│       ├── dependency_injection/   # Внедрение зависимостей
│       ├── entities/               # Сущности и DTO
│       ├── handlers/               # HTTP обработчики
│       ├── middlewares/            # Middleware-ы
│       ├── repositories/           # Работа с данными
│       ├── routes/                 # Инициализация api путей
│       └── services/               # Бизнес-логика
│   ├── pkg/
│       └── utils/                  # Вспомогательные утилиты
│   ├── Dockerfile                  # Конфигурация Docker контейнера для api
│   └── go.mod
├── .env                            # Переменные окружения (нужно создать по образцу см. выше)
├── .gitignore                  
├── docker-compose.yml              # Docker конфигурация
└── README.md                       # README файл
```

## Аутентификация

API использует JWT токены для аутентификации. Для доступа к защищенным endpoint'ам:

1. Выполните вход через `/api/auth/login`
2. Получите JWT токен
3. Добавьте токен в заголовок запроса:
   ```
   Authorization: Bearer <your_jwt_token>
   ```

## Основные endpoint'ы

### Публичные endpoints
- `POST /api/register` - Регистрация пользователя
- `POST /api/login` - Вход в систему
- `POST /api/refresh` - Обновление пары токенов

### Защищенные endpoints (требуют JWT)
- `GET /api/auth/cat/all` - Получить всех котиков
- `POST /api/auth/cat/create` - Создать котика
- `GET /api/auth/cat/:id` - Получить котика по ID
- `PUT /api/auth/cat/mw/:id` - Обновить котика
- `DELETE /api/auth/cat/mw/:id` - Удалить котика

### Фотографии котиков
- `POST /api/auth/cat/mw/:id/photo/add` - Добавить фотографии
- `GET /api/auth/cat/photo/:photoID` - Получить фотографию
- `POST /api/auth/cat/mw/:id/photo/:photoID/primary` - Сделать фото главным
- `DELETE /api/auth/cat/mw/:id/photo/:photoID` - Удалить фотографию

## Базы данных

### PostgreSQL
- **Порт**: 5432
- **База данных**: app_db
- Автоматически создает необходимые таблицы при первом запуске

### Redis
- **Порт**: 6379
- **Используется для**: хранения refresh токенов пользователя

### MinIO
- **Порт**: 9000 (API), 9001 (Console)
- **Хранилище**: фотографии котиков
- Автоматически создает бакеты при первом запуске

## Логирование

Приложение использует Zerolog для структурированного логирования. Логи выводятся в консоль и включают:
- Временные метки
- Уровни логирования (trace, debug, info, warn, error, fatal, panic)
- Статус ответ
- Время выполнения запроса в мс
