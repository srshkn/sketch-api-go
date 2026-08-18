# Sketch API Go

## Содержание
- [О проекте](#о-проекте)
- [Технические требования](#технические-требования)
    - [Go](#go)
    - [Docker](#docker)
    - [Make](#make)
- [Быстрый запуск](#быстрый-запуск)
- [Доступ после старта](#доступ-после-старта)
    - [Старт](#старт)
    - [Доступ после старта](#доступ-после-старта)
- [Структура проекта](#структура-проекта)

Дополнительная информация:

- [Документация по Backend](docs/BACKEND.md)
- [Доступные endpoints](docs/OPENAPI.md)
- [Информация по базе данных](docs/DATABASE.md)
- [Информация по CI/CD.md](docs/CI-CD.md)

## О проекте

Небольшой HTTP API на Go, построенный по подходу contract-first: публичный
контракт описывается в OpenAPI, а модели, интерфейс сервера и встроенная
спецификация генерируются с помощью `oapi-codegen`.

Сейчас сервис предоставляет проверку состояния API и демонстрационный endpoint
регистрации пользователя. После запуска доступны Swagger UI и OpenAPI-схема.

## Технические требования

### Go 
Go [1.25.6](https://go.dev/dl/) или новее. 

### Docker
Docker используется для запуска приложения и PostgreSQL. Для локальной разработки удобнее всего использовать **Docker Desktop**.

- **macOS:** [Инструкции по установке на Mac](https://docs.docker.com/desktop/setup/install/mac-install/)
- **Windows:** [Инструкции по установке на Windows](https://docs.docker.com/desktop/install/windows-install/)
- **Linux:** [Инструкции по установке на Linux](https://docs.docker.com/desktop/install/linux-install/) 

### Make
В проекте используется **Makefile**, который выступает единой точкой входа для всех базовых команд.

- **macOS:** Установлен по умолчанию в составе утилит разработчика. Если команда `make` не найдена, выполните в терминале:
```bash
xcode-select --install
```
- **Windows:** Утилита `make` можно установить, выполнил команду в **PowerShell** от имени администратора:
```powershell
winget install GnuWin32.Make
```
- **Linux (Ubuntu/Debian):** Устанавливается через терминал командой:
```bash
sudo apt update && sudo apt install make
```

## Быстрый запуск

### Старт

Склонируйте репозиторий и перейдите в его директорию:

```bash
git clone <repository-url>
cd sketch-api-go
```

При необходимости обновите зависимости:

```bash
go mod tidy
```

Запустите окружение:
```bash
make compose-dev
```

### Доступ после старта:

- API: `http://localhost:8080`;
- Swagger UI: `http://localhost:8080/docs/`;
- OpenAPI JSON: `http://localhost:8080/openapi.json`.


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
