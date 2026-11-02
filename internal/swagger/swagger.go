package swagger

import (
	"net/http"

	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v5emb"

	v1Generated "sketch-api-go/internal/generated/v1"
)

// Register добавляет маршруты спецификации OpenAPI и интерфейса Swagger UI.
func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /openapi.json", openAPISpec)

	// Встроенные ресурсы Swagger UI позволяют раздавать документацию без отдельных
	// статических файлов; спецификацию интерфейс загружает из /openapi.json.
	swaggerUI := v5emb.NewHandlerWithConfig(swgui.Config{
		Title:       "Sketch API",
		SwaggerJSON: "/openapi.json",
		BasePath:    "/docs/",
		// Значения настроек передаются как JSON: схемы раскрываются на один уровень,
		// а для моделей по умолчанию показываются примеры данных.
		SettingsUI: map[string]string{
			"defaultModelsExpandDepth": "1",
			"defaultModelExpandDepth":  "1",
			"defaultModelRendering":    `"example"`,
		},
	})

	mux.Handle(
		"GET /docs/",
		swaggerUI,
	)
}

// openAPISpec возвращает JSON-спецификацию из сгенерированного пакета API v1.
func openAPISpec(
	w http.ResponseWriter,
	r *http.Request,
) {
	spec, err := v1Generated.GetSpecJSON()
	if err != nil {
		http.Error(
			w,
			"failed to load OpenAPI specification",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json; charset=utf-8",
	)

	_, _ = w.Write(spec)
}
