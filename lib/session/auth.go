package session

import (
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
)

type Auth struct {
	auth *jwtauth.JWTAuth
}

type Config struct {
	Secret string `toml:"secret" env:"JWT_SECRET"`
}

type Claims struct {
	UserID int64 `json:"userId"`
}

func (c *Claims) ToMap() map[string]interface{} {
	return map[string]interface{}{"userId": c.UserID}
}

var (
	AU *Auth

	ErrNotInit       = errors.New("package not initialized")
	ErrInvalidClaims = errors.New("invalid claims")
	ErrClaimUserID   = errors.New("invalid userId in claims")
)

func Setup(confs Config) *Auth {
	AU = &Auth{
		auth: jwtauth.New("HS256", []byte(confs.Secret), nil),
	}
	return AU
}

func (tc *Claims) Valid() error {
	if tc.UserID == 0 {
		return ErrClaimUserID
	}
	return nil
}

func (t *Auth) Authority() *jwtauth.JWTAuth {
	return t.auth
}

func (t *Auth) Encode(claims *Claims) (jwt.Token, string, error) {
	return t.auth.Encode(claims.ToMap())
}

func (t *Auth) Decode(tokenStr string) (jwt.Token, error) {
	return t.auth.Decode(tokenStr)
}
