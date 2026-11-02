package meta

import (
	"context"
	v1GenAPI "sketch-api-go/internal/generated/http/v1"
)

// MetaHandler представляет обработчик служебных запросов API.
type MetaHandler struct{}

// NewHandler создаёт обработчик служебных запросов API.
func NewHandler() *MetaHandler {
	return &MetaHandler{}
}

// GetHealth обрабатывает GET /health и возвращает статус OK.
// Проверяет доступность API без проверки состояния внешних зависимостей.
func GetHealth(ctx context.Context, request v1GenAPI.GetHealthRequestObject) (v1GenAPI.GetHealthResponseObject, error) {
	return v1GenAPI.GetHealth200JSONResponse{
		Status: "OK",
	}, nil
}
