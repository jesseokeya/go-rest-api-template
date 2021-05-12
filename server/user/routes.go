package user

import (
	"github.com/go-chi/chi"
)

func Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/me", func(r chi.Router) {
		r.Get("/", GetSessionUser)
		r.Put("/", UpdateSessionUser)
	})

	r.Route("/{userID}", func(r chi.Router) {
		r.Put("/", UpdateUser)
	})

	return r
}
