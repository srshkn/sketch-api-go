Посмотрел текущий `develop` именно по `internal`. Структура у тебя получается классическая: `handler → service → repository → db`, при этом инфраструктурные вещи вынесены в отдельные пакеты. ([GitHub][1])

Ниже готовый вариант `docs/BACKEND.md`, без описания endpoints.

# Backend

## Содержание

* [Архитектура](#архитектура)
* [Структура `internal`](#структура-internal)
* [App](#app)
* [Config](#config)
* [Handler](#handler)
* [Service](#service)
* [Repository](#repository)
* [Database](#database)
* [Middleware](#middleware)
* [Authentication](#authentication)
* [Password](#password)
* [Cookie](#cookie)
* [Logging](#logging)
* [Swagger](#swagger)
* [Generated](#generated)
* [Обработка ошибок](#обработка-ошибок)
* [Жизненный цикл запроса](#жизненный-цикл-запроса)

## Архитектура

Backend построен вокруг разделения ответственности между слоями приложения.

Основной поток обработки данных:

```text
HTTP Request
     │
     ▼
  Handler
     │
     ▼
  Service
     │
     ▼
 Repository
     │
     ▼
   DB / sqlc
     │
     ▼
 PostgreSQL
```

Дополнительные компоненты подключаются на уровне приложения:

```text
                    ┌─────────────┐
                    │    Config   │
                    └──────┬──────┘
                           │
┌────────────┐      ┌──────▼──────┐      ┌─────────────┐
│ Middleware │ ───► │     App     │ ◄─── │   Swagger   │
└────────────┘      └──────┬──────┘      └─────────────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
           Handler      Service       Token
              │            │
              │            ▼
              │       Repository
              │            │
              └────────────┼────────────┘
                           ▼
                       PostgreSQL
```

Такое разделение позволяет не смешивать HTTP-логику, бизнес-логику и работу с базой данных.

## Структура `internal`

```text
internal/
├── app/                # сборка и запуск приложения
├── config/             # конфигурация приложения
├── cookie/             # работа с authentication cookies
├── db/                 # sqlc-generated database layer
├── errs/               # ошибки приложения
├── generated/v1/       # код, сгенерированный из OpenAPI
├── handler/v1/         # HTTP handlers
├── logging/            # настройка логирования
├── middleware/         # HTTP middleware
├── password/           # хеширование и проверка паролей
├── postgres/           # подключение к PostgreSQL
├── repository/         # интерфейсы доступа к данным
├── service/            # бизнес-логика
├── swagger/            # интеграция Swagger UI
└── token/              # JWT и refresh token
```

Пакет `internal` используется намеренно: код внутри него является внутренней частью приложения и не предназначен для импорта внешними проектами.

---

## App

`internal/app` является composition root приложения.

Здесь создаются и связываются основные зависимости:

* database layer;
* JWT manager;
* cookie manager;
* services;
* handlers;
* middleware;
* Swagger;
* HTTP server.

`app.New()` собирает приложение через dependency injection. Например, `UserService` получает database repository, а `AuthService` — repository и JWT manager.

После создания зависимостей регистрируются HTTP handlers и middleware, после чего создаётся `http.Server`.

Приложение также отвечает за graceful shutdown: при получении сигнала завершения сервер пытается корректно завершить активные соединения в пределах заданного timeout.

Таким образом, `app` не содержит бизнес-логику — его задача заключается в **сборке приложения**.

---

## Config

`internal/config` содержит конфигурационные структуры и логику получения параметров приложения.

Конфигурация разделена по назначению:

```text
config/
├── config.go
├── cookie.go
├── cors.go
├── jwt.go
├── logger.go
├── postgres.go
└── server.go
```

Отдельные конфигурационные типы используются для:

* HTTP server;
* PostgreSQL;
* JWT;
* cookies;
* CORS;
* logger.

Такой подход позволяет не передавать по приложению глобальную конфигурацию целиком, а передавать конкретную конфигурационную секцию нужному компоненту.

---

## Handler

HTTP handlers находятся в `internal/handler/v1`.

Структура:

```text
handler/v1/
├── auth.go
├── handler.go
├── meta.go
├── response.go
└── user.go
```

Handler отвечает за HTTP-уровень:

* получение HTTP request;
* работу с generated OpenAPI-моделями;
* вызов service;
* преобразование результата в HTTP response;
* передачу ошибок в единый механизм формирования ответа.

Handler не должен содержать бизнес-логику или напрямую работать с PostgreSQL.

Общая зависимость выглядит следующим образом:

```text
Handler
   │
   └──► Service
```

Это позволяет независимо тестировать HTTP-слой и бизнес-логику.

---

## Service

`internal/service` содержит бизнес-логику приложения.

```text
service/
├── auth.go
└── user.go
```

### User service

`UserService` отвечает за операции, связанные с пользователями.

Перед созданием пользователя выполняются проверки существования пользователя по email и username. Email нормализуется приведением к lowercase.

При создании пользователя:

```text
Request
  │
  ▼
Validation
  │
  ▼
Password Hash
  │
  ▼
Repository.CreateUser
```

Сервис использует `password` package для получения безопасного хеша пароля и передаёт в repository уже готовый `PasswordHash`.

### Auth service

`AuthService` содержит логику аутентификации и управления токенами.

Он работает через две зависимости:

```go
type AuthService struct {
    repository repository.AuthRepository
    tokenManager token.JWTManager
}
```

При выдаче пары токенов:

```text
User ID
  │
  ├──► Access Token
  │
  └──► Refresh Token
           │
           ▼
       SHA-256 hash
           │
           ▼
       PostgreSQL
```

В базе хранится не сам refresh token, а его хеш. При последующем использовании токена приложение снова вычисляет хеш и ищет соответствующую запись.

При refresh старый refresh token удаляется перед выдачей новой пары токенов.

---

## Repository

`internal/repository` содержит интерфейсы доступа к данным.

```text
repository/
├── auth.go
└── user.go
```

Repository является абстракцией над database layer.

Например, `repository.User` определяет операции, необходимые `UserService`:

```go
type User interface {
    CheckUserExists(...)
    CreateUser(...)
    GetUserByID(...)
}
```

А `repository.Auth` предоставляет операции для аутентификации и refresh tokens.

Важный принцип:

```text
Service → Repository interface → sqlc implementation
```

Service знает только об интерфейсе repository и не зависит от конкретной реализации database layer.

Это упрощает тестирование и позволяет заменить реализацию хранилища без изменения бизнес-логики.

---

## Database

`internal/db` содержит код, сгенерированный `sqlc`, и типы, используемые для работы с PostgreSQL.

```text
db/
├── db.go
├── models.go
├── querier.go
├── refrash_tokens.sql.go
└── users.sql.go
```

SQL-запросы находятся за пределами `internal`, а `sqlc` генерирует Go-код для их выполнения.

Repository использует сгенерированные типы и методы:

```text
SQL
 │
 ▼
sqlc
 │
 ▼
internal/db
 │
 ▼
Repository
 │
 ▼
Service
```

Благодаря этому SQL не находится внутри service или handler.

---

## Middleware

`internal/middleware` содержит HTTP middleware:

```text
middleware/
├── auth.go
├── cors.go
└── logging.go
```

### Authentication middleware

Authentication middleware использует `JWTManager` для проверки access token.

Middleware отвечает за:

* получение токена;
* его валидацию;
* извлечение идентификатора пользователя;
* передачу результата дальше по цепочке request context.

При этом middleware не занимается бизнес-операциями пользователя — это ответственность service layer.

### Logging middleware

Logging middleware интегрирован с `log/slog` и отвечает за логирование HTTP request lifecycle.

### CORS middleware

CORS middleware применяет настройки из `config.CORS` к HTTP server.

---

## Authentication

Для аутентификации используется комбинация:

* short-lived access JWT;
* long-lived refresh token;
* PostgreSQL для хранения refresh token hashes.

Access token подписывается алгоритмом `RS256`.

JWT содержит:

* user ID;
* token type;
* issuer;
* subject;
* issued-at;
* expiration.

Refresh token генерируется криптографически безопасным способом и не является JWT.

Для refresh token используется схема:

```text
Random Token
     │
     ▼
 SHA-256 Hash
     │
     ▼
 PostgreSQL
```

Это позволяет не хранить действующие refresh tokens в открытом виде в базе.

---

## Password

`internal/password` инкапсулирует всю работу с паролями.

Для хеширования используется **Argon2id**.

Текущие параметры:

```text
memory      = 64 MiB
iterations  = 3
parallelism = 4
salt        = 16 bytes
key         = 32 bytes
```

Для каждого пароля генерируется отдельная случайная salt.

Формат сохранённого значения:

```text
$argon2id$v=...$m=...,t=...,p=...$salt$hash
```

Проверка пароля выполняется через повторное вычисление Argon2id и `crypto/subtle.ConstantTimeCompare`.

Таким образом, plaintext password никогда не передаётся в repository и не сохраняется в базе.

---

## Cookie

`internal/cookie` инкапсулирует работу с authentication cookies.

Отдельный пакет позволяет не размазывать настройки cookie по handlers и централизовать:

* имя cookie;
* lifetime;
* path;
* security flags;
* установку и удаление cookie.

Handler работает с абстракцией cookie manager, а не создаёт `http.Cookie` вручную в каждом месте.

---

## Logging

`internal/logging` содержит настройку `log/slog`.

Логгер создаётся один раз при старте приложения и передаётся зависимостям через dependency injection.

Это позволяет:

* использовать единый logger;
* не создавать logger внутри handlers/services;
* централизованно управлять уровнем и форматом логирования.

---

## Swagger

`internal/swagger` отвечает за интеграцию Swagger UI с HTTP server.

Swagger использует OpenAPI specification, которая является исходным контрактом API.

Общий поток:

```text
api/v1/openapi.yml
        │
        ▼
   oapi-codegen
        │
        ▼
internal/generated/v1
        │
        ├── models.gen.go
        ├── server.gen.go
        └── spec.gen.go
```

Swagger UI использует ту же спецификацию, что и generated API layer, поэтому документация и реализация должны оставаться синхронизированными через OpenAPI contract.

Описание самого API вынесено в отдельную документацию: [`OPENAPI.md`](./OPENAPI.md).

---

## Generated

`internal/generated/v1` содержит код, автоматически сгенерированный `oapi-codegen`:

```text
generated/v1/
├── models.gen.go
├── server.gen.go
└── spec.gen.go
```

Generated code не редактируется вручную.

Источником истины является:

```text
api/v1/openapi.yml
```

После изменения OpenAPI-контракта generated code должен быть пересоздан.

Это позволяет использовать contract-first подход:

```text
OpenAPI Contract
       │
       ▼
  oapi-codegen
       │
       ▼
Generated Go API
       │
       ▼
Handlers
```

---

## Обработка ошибок

Ошибки разделены по уровням ответственности.

Service может возвращать ошибки бизнес-логики, например:

* пользователь уже существует;
* неверный пароль;
* refresh token истёк.

Repository возвращает ошибки, связанные с хранилищем данных.

Handler не должен самостоятельно реализовывать бизнес-логику обработки этих ошибок. Его задача — передать результат в общий механизм формирования HTTP response.

Для типизированных ошибок используются отдельные определения в `internal/errs`.

---

## Жизненный цикл запроса

В общем случае запрос проходит через следующие этапы:

```text
HTTP Request
     │
     ▼
CORS Middleware
     │
     ▼
Logging Middleware
     │
     ▼
Auth Middleware
     │
     ▼
Generated Server
     │
     ▼
Handler
     │
     ▼
Service
     │
     ▼
Repository
     │
     ▼
sqlc / PostgreSQL
     │
     ▼
Repository
     │
     ▼
Service
     │
     ▼
Handler
     │
     ▼
HTTP Response
```

Не каждый запрос использует все перечисленные компоненты. Например, операции, не требующие аутентификации, не используют authentication middleware на уровне бизнес-операции.

Главная идея архитектуры — **каждый слой отвечает только за свою область**:

| Слой       | Ответственность             |
| ---------- | --------------------------- |
| Handler    | HTTP                        |
| Service    | Бизнес-логика               |
| Repository | Абстракция доступа к данным |
| DB         | SQL и PostgreSQL            |
| Middleware | Сквозные HTTP-механизмы     |
| Token      | JWT и refresh tokens        |
| Password   | Хеширование паролей         |
| Config     | Конфигурация                |
| App        | Сборка приложения           |
| Generated  | Код из OpenAPI-контракта    |

Такое разделение является базовым шаблоном проекта и позволяет добавлять новые домены без необходимости изменять существующую архитектуру.

[1]: https://github.com/srshkn/sketch-api-go/tree/develop/internal "sketch-api-go/internal at develop · srshkn/sketch-api-go · GitHub"
