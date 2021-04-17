package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

// Status customizes http response status
type Status struct {
	Code int
	Text string
}

func renderResponder(w http.ResponseWriter, r *http.Request, v interface{}) {
	if v, ok := v.(*ApiError); ok {
		LogEntrySetFields(r, map[string]interface{}{
			"error":    v.Error(),
			"location": v.Location,
		})
	}
	if v, ok := v.(Status); ok {
		w.WriteHeader(v.Code)
		LogEntrySetField(r, "status", v.Code)
	}
	render.DefaultResponder(w, r, v)
}

func Render(w http.ResponseWriter, r *http.Request, v render.Renderer) {
	if err := render.Render(w, r, v); err != nil {
		log.Error().Timestamp().
			Str("error", err.Error()).
			Send()
	}
}

func init() {
	// inject and override defaults in the render package
	render.Respond = renderResponder
}
