package swagger

import (
	"net/http"

	"github.com/swaggest/swgui/v5emb"

	"sketch-api-go/internal/generated"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /openapi.json", openAPISpec)

	mux.Handle(
		"GET /docs/",
		v5emb.New(
			"Sketch API",
			"/openapi.json",
			"/docs/",
		),
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
