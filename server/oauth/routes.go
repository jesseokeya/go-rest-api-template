package oauth

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/jesseokeya/go-rest-api-template/server/auth"
)

func GetPing(w http.ResponseWriter, r *http.Request) {
	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		fmt.Fprintf(w, "Permission: %v", err)
		return
	}
	userId, _ := token.Get("userId")
	fmt.Fprintf(w, "User ID: %v", userId)
}

func (o *OAuth) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/signin", o.Signin)
	r.Post("/signup", auth.Signup)
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(o.Authority()))

		// Handle valid / invalid tokens.
		r.Use(auth.Authenticator)

		r.Get("/ping", GetPing)
		r.Get("/authorize", o.Authorize)
	})
	r.Get("/me", o.GetMe)
	r.Post("/token", o.GetToken)

	return r
}
