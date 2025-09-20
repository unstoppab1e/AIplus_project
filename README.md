# AIPlus Employee management service

Веб-сервис для управления данными сотрудников на Go с PostgreSQL.

## Функциональность

- ✅ Создание сотрудника (ФИО, телефон, город)
- ✅ Получение сотрудника по ID
- ✅ Получение списка всех сотрудников
- ✅ Обновление данных сотрудника
- ✅ Удаление сотрудника
- ✅ Валидация данных
- ✅ Миграции базы данных
- ✅ Docker контейнеризация
- ✅ Unit тесты

## Быстрый старт

### С Docker Compose (рекомендуется)

1. Клонируйте репозиторий
2. Запустите приложение:
```bash
docker-compose up --build
```

3. Сервис будет доступен по адресу: http://localhost:8080

### Локальная разработка

1. Установите зависимости:
```bash
go mod tidy
```

2. Запустите PostgreSQL:
```bash
docker run --name postgres -e POSTGRES_DB=employees -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:15
```

3. Установите переменные окружения:
```bash
export DATABASE_URL="postgres://postgres:password@localhost:5432/employees?sslmode=disable"
export PORT="8080"
```

4. Запустите приложение:
```bash
go run cmd/server/main.go
```

## API Endpoints

### Health Check
```bash
GET /health
```

### Создать сотрудника
```bash
POST /api/v1/employees
Content-Type: application/json

{
  "full_name": "Айгүл Серікқызы",
  "phone": "+7 701 987 65 43",
  "city": "Алматы"
}
```

### Получить всех сотрудников
```bash
GET /api/v1/employees
```

### Получить сотрудника по ID
```bash
GET /api/v1/employees/{id}
```

### Обновить сотрудника
```bash
PUT /api/v1/employees/{id}
Content-Type: application/json

'{
  "full_name": "Әлихан Әбдірахман",
  "phone": "+7 777 123 45 67",
  "city": "Астана"
}
```

### Удалить сотрудника
```bash
DELETE /api/v1/employees/{id}
```

## Примеры использования с curl

### Создание сотрудника
```bash
curl -X POST http://localhost:8080/api/v1/employees \
  -H "Content-Type: application/json" \
  -d '{
  "full_name": "Айгүл Серікқызы",
  "phone": "+7 701 987 65 43",
  "city": "Алматы"
}'
```

### Получение всех сотрудников
```bash
curl http://localhost:8080/api/v1/employees
```

### Получение сотрудника по ID
```bash
curl http://localhost:8080/api/v1/employees/1
```

### Обновление сотрудника
```bash
curl -X PUT http://localhost:8080/api/v1/employees/1 \
  -H "Content-Type: application/json" \
  -d '{
  "full_name": "Әлихан Әбдірахман",
  "phone": "+7 777 123 45 67",
  "city": "Астана"
}'
```

### Удаление сотрудника
```bash
curl -X DELETE http://localhost:8080/api/v1/employees/1
```

## Тестирование

### Запуск всех тестов
```bash
go test ./tests/... -v
```

### Запуск тестов с покрытием
```bash
go test ./tests/... -v -cover
```

### Запуск конкретного теста
```bash
go test ./tests/... -v -run TestEmployeeHandler_CreateEmployee
```

## Структура проекта

```
employee-service/
├── cmd/
│   └── server/
│       └── main.go              # Точка входа в приложение
├── internal/
│   ├── config/
│   │   └── config.go            # Конфигурация приложения
│   ├── handler/
│   │   └── employee.go          # HTTP обработчики
│   ├── model/
│   │   └── employee.go          # Модели данных
│   ├── repository/
│   │   └── employee.go          # Слой доступа к данным
│   └── service/
│       └── employee.go          # Бизнес-логика
├── migrations/
│   └── 001_create_employees_table.down.sql  # SQL миграции
│   └── 001_create_employees_table.up.sql  # SQL миграции
├── tests/
│   ├── handler_test.go          # Тесты HTTP слоя
│   └── repository_test.go       # Тесты репозитория
├── docker-compose.yml           # Docker Compose конфигурация
├── Dockerfile                   # Docker образ приложения
├── go.mod                       # Go модуль
├── go.sum                       # Хеши зависимостей
└── README.md                    # Документация
```

## Переменные окружения

- `DATABASE_URL` - строка подключения к PostgreSQL (обязательно)
- `PORT` - порт для HTTP сервера (по умолчанию 8080)

## Troubleshooting

### Проблемы с запуском

**Сервис не запускается**
```bash
# Проверьте логи
docker-compose logs

# Пересоберите образы
docker-compose build --no-cache
docker-compose up
```

**Ошибка подключения к БД**
```bash
# Убедитесь что PostgreSQL запустился
docker-compose logs db

# Проверьте порты
docker-compose ps
```

**Проблемы с миграциями**
```bash
# Очистите данные и пересоздайте
docker-compose down -v
docker-compose up --build
```

### Быстрое решение проблем

```bash
# Полный сброс
docker-compose down -v
docker system prune -f
docker-compose up --build

# Проверка работы
./test-docker.sh
```
