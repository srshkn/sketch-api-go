# Backend

## Обзор

Sketch API Go — небольшой HTTP-сервис на стандартном `net/http`. API
разрабатывается по подходу contract-first: источником публичного контракта
является `api/openapi.yml`, а `oapi-codegen` создаёт из него:

- Go-модели запросов и ответов;
- интерфейс сервера и регистрацию маршрутов;
- встроенную OpenAPI-спецификацию для Swagger UI.

Точка входа `cmd/server/main.go` создаёт `http.ServeMux`, подключает API-handlers
и Swagger, после чего запускает сервер на порту `8080`.

Базовый URL для локальной разработки: `http://localhost:8080`.

## Технологии

- Go и стандартный пакет `net/http`;
- OpenAPI 3.0.3;
- `oapi-codegen` для генерации серверного кода;
- `swaggest/swgui` для Swagger UI;
- `httptest` для endpoint-тестов;
- multi-stage Docker-сборка и distroless runtime-образ.

## Маршрутизация запроса

```text
HTTP request
    ↓
http.ServeMux
    ↓
generated.ServerInterfaceWrapper
    ↓
internal/handlers.Handler
    ↓
HTTP response
```

Сгенерированный слой связывает пути из OpenAPI с реализацией
`internal/handlers`. Тип `Handler` реализует интерфейс
`generated.ServerInterface`, поэтому отсутствие обязательного handler будет
обнаружено во время компиляции.

## Endpoints

| Метод | Путь | Назначение | Успешный статус |
| --- | --- | --- | --- |
| `GET` | `/health` | Проверка доступности сервиса | `200 OK` |
| `POST` | `/user/register` | Демонстрационная регистрация пользователя | `201 Created` |
| `GET` | `/docs/` | Swagger UI | `200 OK` |
| `GET` | `/openapi.json` | OpenAPI-контракт в JSON | `200 OK` |

### GET /health

Проверяет, что процесс API запущен и принимает HTTP-запросы.

```bash
curl http://localhost:8080/health
```

Ответ `200 OK` с заголовком `Content-Type: text/plain; charset=utf-8`:

```text
OK
```

### POST /user/register

Принимает имя и пароль пользователя. Имя очищается от пробелов в начале и
конце строки.

Тело запроса:

| Поле | Тип | Обязательное | Описание |
| --- | --- | --- | --- |
| `name` | string | да | Имя пользователя |
| `password` | string | да | Пароль пользователя |

Успешный ответ имеет статус `201 Created` и заголовок
`Content-Type: application/json`.

| Поле ответа | Тип | Текущее значение |
| --- | --- | --- |
| `id` | integer | `1` |
| `name` | string | Очищенное имя из запроса |

Текущая реализация является демонстрационной: пользователь не сохраняется,
пароль не хешируется, а в успешном ответе всегда возвращается `id: 1`.

#### Валидация

Handler возвращает `400 Bad Request`, если:

- тело пустое или содержит некорректный JSON;
- запрос содержит неизвестное поле;
- `name` отсутствует, пуст или состоит только из пробелов;
- `password` отсутствует или является пустой строкой.

Ответ об ошибке содержит объект `error` с кодом `INVALID_REQUEST` и сообщением,
например `name must not be empty`.

В `api/openapi.yml` для `name` также указан `minLength: 5`, а для `password` —
`minLength: 8`. На текущем этапе эти ограничения описаны в контракте, но не
проверяются handler вручную. До добавления middleware-валидации фактическим
поведением являются проверки, перечисленные выше.

## Формат ошибок

Ошибки API, предусмотренные `ErrorResponse`, имеют поля:

| Поле | Тип | Назначение |
| --- | --- | --- |
| `error.code` | string | Машиночитаемый код ошибки |
| `error.message` | string | Понятное человеку описание |

OpenAPI-контракт содержит общий набор кодов для дальнейшего развития API:
`INVALID_REQUEST`, `UNAUTHORIZED`, `NOT_FOUND`, `ROOM_NOT_FOUND`,
`SLOT_NOT_FOUND`, `SLOT_ALREADY_BOOKED`, `BOOKING_NOT_FOUND`, `FORBIDDEN`,
`SCHEDULE_EXISTS` и `INTERNAL_ERROR`. Текущие handlers используют
`INVALID_REQUEST`.

## Структура backend

```text
api/
├── openapi.yml              # источник API-контракта
└── configs/
    ├── models.yml           # генерация моделей
    ├── server.yml           # генерация интерфейса и роутера
    └── spec.yml             # встраивание OpenAPI-спецификации

cmd/server/
└── main.go                  # сборка зависимостей и запуск HTTP-сервера

internal/
├── generated/
│   ├── models.gen.go        # модели API
│   ├── server.gen.go        # ServerInterface и маршрутизация
│   └── spec.gen.go          # встроенная спецификация
├── handlers/
│   ├── handlers.go          # бизнес-логика HTTP-handlers
│   └── handlers_test.go     # endpoint-тесты
└── swagger/
    └── swagger.go           # Swagger UI и /openapi.json
```

Пакет `internal` ограничивает использование backend-компонентов текущим
Go-модулем.

## Изменение API

При добавлении или изменении endpoint:

1. Обновите контракт в `api/openapi.yml`.
2. Перегенерируйте модели, интерфейс сервера и встроенную спецификацию:

   ```bash
   make gen
   ```

3. Реализуйте новый метод интерфейса в `internal/handlers`.
4. Добавьте endpoint-тесты в `internal/handlers`.
5. Запустите форматирование и проверки:

   ```bash
   gofmt -w .
   go vet ./...
   go test ./...
   go build ./...
   ```

Сгенерированные файлы в `internal/generated` не редактируются вручную.

Если `make` недоступен, генерацию можно выполнить напрямую:

```bash
go tool oapi-codegen -config ./api/configs/server.yml ./api/openapi.yml
go tool oapi-codegen -config ./api/configs/models.yml ./api/openapi.yml
go tool oapi-codegen -config ./api/configs/spec.yml ./api/openapi.yml
```

## Локальный запуск и тесты

Запустить сервер без повторной генерации:

```bash
go run ./cmd/server
```

Перегенерировать код и запустить сервер:

```bash
make run
```

Запустить весь набор тестов:

```bash
go test ./...
```

Запустить endpoint-тесты с подробным выводом:

```bash
go test -v ./internal/handlers
```

## Docker

Собрать и запустить образ:

```bash
docker build -t sketch-api-go .
docker run --rm --name sketch-api-go -p 8080:8080 sketch-api-go
```

Dockerfile компилирует статический Linux-бинарный файл в builder-образе, затем
копирует его в минимальный distroless-образ и запускает от непривилегированного
пользователя.
