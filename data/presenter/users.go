package presenter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clout-jam/api/data"
	"github.com/clout-jam/api/lib/session"
)

// User holds a pointer to the postgres user data structure
type User struct {
	*data.User
	ID         int64       `json:"id"`
	Name       string      `json:"name"`
	HasProfile bool        `json:"hasProfile"`
	Etc        interface{} `json:"etc"`
	CreatedAt  interface{} `json:"createdAt"`
	UpdatedAt  interface{} `json:"updatedAt"`
}

// Homework holds refernce to user lessons
type Homework struct {
	*data.UserLesson
	*data.Lesson
	*User
}

func (u *User) Render(w http.ResponseWriter, r *http.Request) error {
	u.HasProfile = u.User.Etc.HasProfile
	return nil
}

// NewUser creates a new user in postgres
func NewUser(ctx context.Context, user *data.User) *User {
	presented := &User{
		User: user,
	}

	presented.Name = user.FirstName
	if user.LastName != "" {
		presented.Name = fmt.Sprintf("%s %s", presented.Name, user.LastName)
	}

	return presented
}

type AuthUser struct {
	*User
	JWT   string `json:"jwt"`
	Token string `json:"token"`
}

type AuthIntercom struct {
	HashID string `json:"hashId"`
	Hash   string `json:"hash"`
}

func (u *AuthUser) Render(w http.ResponseWriter, r *http.Request) error {
	claims := &session.Claims{UserID: u.ID}
	_, u.JWT, _ = session.AU.Encode(claims)
	// backward compat:
	u.Token = u.JWT
	return u.User.Render(w, r)
}

func NewAuthUser(ctx context.Context, user *data.User) *AuthUser {
	authUser := &AuthUser{
		User: NewUser(ctx, user),
	}
	return authUser
}
