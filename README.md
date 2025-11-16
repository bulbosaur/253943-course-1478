# Order Service

gRPC-сервис для управления заказами. Поддерживает создание, получение, обновление, удаление и список заказов.

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

Создайте файл `.env` на основе шаблона в разделе [конфигурация](#конфигурация)

Можете переименовать существующий `env.example` по адресу `.\config\env.example` в `.env`.

### 3. Запуск с помощью Docker Compose

```
docker-compose up --build
```

Это автоматически:

- Соберёт образ приложения.
- Запустит контейнеры приложения, PostgreSQL и Redis.
- Выполнит миграции базы данных при запуске.

## Конфигурация

### Пример `.env`

```
ENV_LOGLEVEL=info           # Уровень логирования: debug, info

GRPC_PORT=50051             # Порт gRPC-сервера
GRPC_HOST=localhost         # Хост gRPC-сервера

HTTP_PORT=8080              # Порт HTTP-сервера
HTTP_HOST=localhost         # Хост HTTP-сервера
HTTP_TIMEOUT=30s            # Таймаут для HTTP-запросов

POSTGRESQL_USER=postgres    # Имя пользователя PostgreSQL
POSTGRES_PASSWORD=postgres  # Пароль PostgreSQL
POSTGRES_HOST=localhost     # Адрес PostgreSQL
POSTGRES_PORT=5432          # Порт PostgreSQL
POSTGRES_DB=orders          # Имя базы данных

REDIS_HOST=localhost        # Адрес Redis
REDIS_PORT=6379             # Порт Redis
REDIS_PASSWORD=redis        # Пароль Redis (пусто — если без аутентификации)
REDIS_DB=0                  # Номер логической БД в Redis (0–15)
REDIS_MAX_RETRIES=5         # Макс. число повторных попыток при ошибках
REDIS_DIAL_TIMEOUT=10s      # Таймаут установки соединения с Redis
REDIS_TIMEOUT=5s            # Общий таймаут на чтение/запись в Redis

```

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

При запуске REST-сервера на :8080, будут доступны следующие запросы:

| Метод GRPC    | HTTP метод | HTTP ручка                                                |
| ------------- | ---------- | --------------------------------------------------------- |
| `CreateOrder` | `POST`     | `http://localhost:8080/api.OrderService/CreateOrder`      |
| `GetOrder`    | `GET `     | `http://localhost:8080/api.OrderService/GetOrder/{id}`    |
| `UpdateOrder` | `POST`     | `http://localhost:8080/api.OrderService/UpdateOrder`      |
| `DeleteOrder` | `DELETE`   | `http://localhost:8080/api.OrderService/DeleteOrder/{id}` |
| `ListOrders`  | `GET`      | `http://localhost:8080/api.OrderService/ListOrders`       |

Для отладки можно использовать:

- Postman
