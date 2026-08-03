package swagger

import (
	"net/http"

	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v5emb"

	"sketch-api-go/internal/generated"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /openapi.json", openAPISpec)

	swaggerUI := v5emb.NewHandlerWithConfig(swgui.Config{
		Title:       "Sketch API",
		SwaggerJSON: "/openapi.json",
		BasePath:    "/docs/",
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

func openAPISpec(
	w http.ResponseWriter,
	r *http.Request,
) {
	spec, err := generated.GetSpecJSON()
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
