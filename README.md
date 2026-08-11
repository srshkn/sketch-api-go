# Sketch API Go

Небольшой HTTP API на Go, построенный по подходу contract-first: публичный
контракт описывается в OpenAPI, а модели, интерфейс сервера и встроенная
спецификация генерируются с помощью `oapi-codegen`.

Сейчас сервис предоставляет проверку состояния API и демонстрационный endpoint
регистрации пользователя. После запуска доступны Swagger UI и OpenAPI-схема.

Подробное описание backend и API находится в
[docs/BACKEND.md](docs/BACKEND.md).

## Быстрый запуск

Требования:

- Go 1.25.6 или новее;
- Git;
- опционально: `make` и Docker.

Склонируйте репозиторий и перейдите в его директорию:

```bash
git clone <repository-url>
cd sketch-api-go
```

Загрузите зависимости и запустите сервер:

```bash
go mod download
go run ./cmd/server
```

Сервис запустится на `http://localhost:8080`.

Проверить его работу можно запросом:

```bash
curl http://localhost:8080/health
```

Ожидаемый ответ:

```text
OK
```

После запуска доступны:

- API: `http://localhost:8080`;
- Swagger UI: `http://localhost:8080/docs/`;
- OpenAPI JSON: `http://localhost:8080/openapi.json`.

## Запуск через Make

```bash
make run
```

Команда сначала обновит сгенерированный Go-код из `api/openapi.yml`, а затем
запустит сервер. Список остальных команд:

```bash
make help
```

## Запуск в Docker

```bash
docker build -t sketch-api-go .
docker run --rm --name sketch-api-go -p 8080:8080 sketch-api-go
```

Либо с помощью Make:

```bash
make docker-build
make docker-run
```

## Тестирование

Запустить все проверки пакетов:

```bash
go test ./...
```

Запустить только тесты HTTP-handlers:

```bash
go test ./internal/handlers
```

## Структура проекта

```text
sketch-api-go
│
├── .github/workflows/       # CI/CD
│
├── api/v1
│   ├── configs/             # конфигурация генератора oapi-codegen
│   └── openapi.yml          # исходный OpenAPI-контракт
│
├── cmd/
│   └── server/
│       └── main.go          # точка входа приложения
│
├── db/                      # SQL-миграции и запросы к PostgreSQL
│
├── docs/                    # доп. документация проекта
│
├── infra/                   # инфраструктурный код и конфигурация
│
├── internal/                # внутренний код приложения
│
├── secrets/                 # локальные секреты для разработки
│
├── .dockerignore            # файлы, исключаемые из Docker build context 
├── .env.example             # пример переменных окружения
├── .gitignore               # файлы, исключаемые из Git
├── .version                 # версия приложения
├── compose.yml              # основная Docker Compose конфигурация
├── Dockerfile               # сборка production-образа
├── Makefile                 # команды запуска
├── go.mod                   # модуль и зависимости Go
├── go.sum                   # контрольные суммы зависимостей
├── sqlc.yml                 # конфигурация генерации sqlc
├── LICENSE                  # лицензия проекта
└── README.md                # описание и документация проекта
```

Файлы в `internal/generated` не следует редактировать вручную: они
перезаписываются при выполнении `make gen`.
