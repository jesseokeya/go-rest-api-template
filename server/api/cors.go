package api

import (
	"net/http"

	"github.com/go-chi/cors"
)

func allowOriginFunc(r *http.Request, origin string) bool {
	if origin == "*" {
		return true
	}
	return true
}

// Cors middleware
func Cors() *cors.Cors {
	cors := cors.New(cors.Options{
		AllowOriginFunc: allowOriginFunc,
		AllowedMethods:  []string{"GET", "POST", "PUT", "UPDATE", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		AllowCredentials: true,
		MaxAge:           300,
	})
	return cors
}
