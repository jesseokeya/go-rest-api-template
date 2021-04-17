package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/jesseokeya/go-rest-api-template/data"
	"github.com/jesseokeya/go-rest-api-template/server/api"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/upper/db/v4"
)

var (
	// SessionCtxKey is the context.Context key to store the request context.
	SessionUserCtxKey = &api.ContextKey{Name: "Session.User"}
	SessionRoleCtxKey = &api.ContextKey{Name: "Session.Role"}

	ErrNoToken          = errors.New("no token context")
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidClaimUser = errors.New("invalid claim userId")
	ErrSessionUser      = errors.New("invalid session user")
)

// Authenticator enforces the validity of the jwt token.
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			api.Render(w, r, api.ErrUnauthorized(ErrNoToken))
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			api.Render(w, r, api.ErrUnauthorized(ErrInvalidToken))
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

func SessionCtx(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token, claims, err := jwtauth.FromContext(ctx)
		if err != nil {
			api.Render(w, r, api.ErrUnauthorized(ErrNoToken))
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			api.Render(w, r, api.ErrUnauthorized(ErrInvalidToken))
			return
		}

		rawUserID, ok := claims["userId"].(float64)
		if !ok {
			api.Render(w, r, api.ErrUnauthorized(ErrInvalidClaimUser))
			return
		}
		userID := int64(rawUserID)

		user, err := data.DB.User.FindByID(userID)
		if err != nil {
			if err == db.ErrNoMoreRows {
				api.Render(w, r, api.ErrUnauthorized(ErrSessionUser))
				return
			}
			api.Render(w, r, api.ErrServiceUnavailable(err))
			return
		}

		ctx = context.WithValue(ctx, SessionRoleCtxKey, user.Role)
		ctx = context.WithValue(ctx, SessionUserCtxKey, user)

		api.LogEntrySetField(r, "sessionId", user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
