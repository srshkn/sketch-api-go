# Sketch API on Go

API предоставляет базовые возможности для регистрации пользователей, аутентификации, управления сессией и получения информации о текущем пользователе.

## Base URL

Все API-эндпоинты находятся под префиксом:

```text
/api/v1
```

Например:

```text
GET /api/v1/health
```

## Аутентификация

Для защищённых эндпоинтов используется **Bearer JWT**.

Заголовок запроса:

```http
Authorization: Bearer <access_token>
```

Access token выдаётся после успешной авторизации через:

```http
POST /api/v1/auth/login
```

Защищённые эндпоинты:

* `GET /user/me`
* `POST /auth/logout`

---

## Meta

### `GET /health`

Проверяет доступность API.

**Ответ `200 OK`:**

```text
OK
```

**Ответ `500 Internal Server Error`:**

```json
{
  "message": "internal server error"
}
```

---

# User

## `POST /user/register`

Регистрирует нового пользователя.

### Request

```json
{
  "username": "john",
  "email": "john@example.com",
  "password": "password123",
  "confirmation": "password123"
}
```

| Поле           | Тип    | Ограничения                            |
| -------------- | ------ | -------------------------------------- |
| `username`     | string | минимум 1 символ                       |
| `email`        | string | корректный email, максимум 254 символа |
| `password`     | string | 8–128 символов                         |
| `confirmation` | string | 8–128 символов                         |

### Response `201 Created`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "john"
}
```

### Возможные ошибки

* `400 Bad Request` — некорректные данные запроса.
* `409 Conflict` — пользователь уже существует.
* `500 Internal Server Error` — внутренняя ошибка сервера.

---

## `GET /user/me`

Возвращает информацию о текущем авторизованном пользователе.

Требует `access_token`.

### Request

```http
Authorization: Bearer <access_token>
```

### Response `200 OK`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "john"
}
```

### Возможные ошибки

* `401 Unauthorized` — отсутствует или недействителен access token.
* `500 Internal Server Error` — внутренняя ошибка сервера.

---

# Auth

## `POST /auth/login`

Аутентифицирует пользователя и выдаёт access token.

### Request

```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

### Response `200 OK`

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

Access token необходимо передавать в последующих защищённых запросах:

```http
Authorization: Bearer <access_token>
```

### Возможные ошибки

* `400 Bad Request` — некорректные данные запроса.
* `409 Conflict` — пользователь не существует или данные аутентификации некорректны.
* `500 Internal Server Error` — внутренняя ошибка сервера.

---

## `POST /auth/logout`

Завершает текущую авторизованную сессию.

Требует `access_token`.

### Request

```http
Authorization: Bearer <access_token>
```

### Response `204 No Content`

Тело ответа отсутствует.

### Возможные ошибки

* `401 Unauthorized` — отсутствует или недействителен access token.
* `500 Internal Server Error` — внутренняя ошибка сервера.

---

## `POST /auth/refresh`

Обновляет access token.

> Механизм refresh token находится в разработке. В текущей OpenAPI-схеме `refresh_token` не возвращается в `TokensResponse`.

### Response `200 OK`

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Возможные ошибки

* `400 Bad Request` — некорректный запрос.
* `409 Conflict` — пользователь не существует.
* `500 Internal Server Error` — внутренняя ошибка сервера.

---

# Error Handling

Ошибки API, связанные с клиентским запросом, возвращаются в формате:

```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "invalid request"
  }
}
```

Поддерживаемые коды:

| Code              | Назначение                |
| ----------------- | ------------------------- |
| `INVALID_REQUEST` | Некорректный запрос       |
| `UNAUTHORIZED`    | Требуется авторизация     |
| `NOT_FOUND`       | Ресурс не найден          |
| `FORBIDDEN`       | Доступ запрещён           |
| `SCHEDULE_EXISTS` | Ресурс уже существует     |
| `INTERNAL_ERROR`  | Внутренняя ошибка сервера |

Внутренние ошибки сервера могут возвращаться в формате:

```json
{
  "message": "internal server error"
}
```

---

# API Overview

| Метод  | Endpoint         | Auth   | Назначение              |
| ------ | ---------------- | ------ | ----------------------- |
| `GET`  | `/health`        | —      | Проверка API            |
| `POST` | `/user/register` | —      | Регистрация             |
| `GET`  | `/user/me`       | Bearer | Текущий пользователь    |
| `POST` | `/auth/login`    | —      | Авторизация             |
| `POST` | `/auth/logout`   | Bearer | Выход                   |
| `POST` | `/auth/refresh`  | —      | Обновление access token |

## OpenAPI

Спецификация API описывается в формате **OpenAPI 3.0.3** и используется для генерации HTTP-интерфейсов и Swagger UI.