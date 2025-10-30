# Order Service

gRPC-сервис для управления заказами. Поддерживает создание, получение, обновление, удаление и список заказов.

---

## Структура проекта

```
.
├── api/                 # .proto-файлы
├── cmd/                 # Точка входа (main.go)
├── config/              # Конфигурация (.env, config.yaml)
├── internal/            # Приватная логика приложения
│   ├── storage/         # Хранилище заказов (in-memory)
│   └── transport/gRPC/  # gRPC-хендлеры и сервер
├── logger/              # Логирование (на базе zap)
├── pkg/api/test/        # Сгенерированный gRPC-код
├── .golangci.yml        # Конфигурация линтера
├── go.mod               # Зависимости Go
└── Makefile             # Скрипты сборки и разработки
```

---

## Требования

- **Go** 1.21+
- **protoc** (Protocol Buffers Compiler) ≥ 3.15
- **make** (опционально, но рекомендуется)

---

## Настройка и запуск

### 1. Клонирование и подготовка

```
git clone https://gitlab.crja72.ru/golang/2025/spring/course/students/253943-Sofiytula71-gmail.com-course-1478
cd 253943-Sofiytula71-gmail.com-course-1478
```

### 2. Настройка окружения

Создайте файл `.env` на основе шаблона:

```
ENV_LOGLEVEL=debug      # Уровень логирования

GRPC_PORT=50051         # Порт запуска gRPC сервера
GRPC_HOST=localhost     # Хост запуска gRPC сервера
```

### 3. Генерация gRPC-кода

```
make generate
```

> Эта команда прочитает `api/order.proto` и сгенерирует Go-файлы в `pkg/api/test/`.

### 4. Сборка и запуск

```
make run
```

Сервер запустится на адресе, указанном в конфигурации (по умолчанию `localhost:50051`).

---

## Конфигурация

Приложение использует **гибридную конфигурацию**:

- Базовые настройки — из `config/config.yaml`
- Переменные окружения — из `config/.env` (или системных переменных)

### Пример `.env`

```
ENV_LOGLEVEL=debug

GRPC_PORT=50051
GRPC_HOST=localhost

HTTP_PORT=8080
HTTP_HOST=localhost
HTTP_TIMEOUT=30s

```

### Описание переменных

| Переменная     | По умолчанию | Описание                                               |
| -------------- | ------------ | ------------------------------------------------------ |
| `ENV_LOGLEVEL` | `info`       | Уровень логирования (`debug`, `info`, `warn`, `error`) |
| `GRPC_HOST`    | `0.0.0.0`    | Хост, на котором слушает gRPC-сервер                   |
| `GRPC_PORT`    | `50051`      | Порт gRPC-сервера                                      |
| `GRPC_HOST`    | `0.0.0.0`    | Хост, на котором слушает HTTP-сервер                   |
| `HTTP_PORT`    | `8080`       | Порт HTTP-сервера                                      |
| `HTTP_timeout` | `30s`        | таймаут HTTP-сервера                                   |

---

## Доступные команды Make

```
make build     # Собрать бинарник
make run       # Собрать и запустить сервер
make generate  # Пересоздать gRPC-код из .proto
make lint      # Проверить код линтером
make test      # Запустить тесты
make clean     # Удалить бинарник
make help      # Показать справку
```

---

## Тестирование

В настоящий момент проект не содержит unit-тестов, но вы можете добавить их в соответствующие пакеты (`internal/storage`, `internal/transport/gRPC` и т.д.).

Запуск тестов:

```
make test
```

---

## gRPC API

Сервис реализует следующие методы:

| Метод         | Запрос               | Ответ                 |
| ------------- | -------------------- | --------------------- |
| `CreateOrder` | `CreateOrderRequest` | `CreateOrderResponse` |
| `GetOrder`    | `GetOrderRequest`    | `GetOrderResponse`    |
| `UpdateOrder` | `UpdateOrderRequest` | `UpdateOrderResponse` |
| `DeleteOrder` | `DeleteOrderRequest` | `DeleteOrderResponse` |
| `ListOrders`  | `ListOrdersRequest`  | `ListOrdersResponse`  |

Для отладки можно использовать:

- Postman (с поддержкой gRPC)

> Сервер включает **gRPC Reflection**, поэтому клиенты могут автоматически обнаруживать методы.

## REST HTTP API

При запуске REST-сервера на :8080, будут доступны следующие POST-запросы:

| Метод GRPC    | HTTP метод | HTTP ручка                                           |
| ------------- | ---------- | ---------------------------------------------------- |
| `CreateOrder` | `POST`     | `http://localhost:8080/api.OrderService/CreateOrder` |
| `GetOrder`    | `POST`     | `http://localhost:8080/api.OrderService/GetOrder`    |
| `UpdateOrder` | `POST`     | `http://localhost:8080/api.OrderService/UpdateOrder` |
| `DeleteOrder` | `POST`     | `http://localhost:8080/api.OrderService/DeleteOrder` |
| `ListOrders`  | `POST`     | `http://localhost:8080/api.OrderService/ListOrders`  |

Для отладки можно использовать:

- Postman
