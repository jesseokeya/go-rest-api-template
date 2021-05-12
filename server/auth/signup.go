package auth

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/goware/emailx"
	"github.com/jesseokeya/go-rest-api-template/data"
	"github.com/jesseokeya/go-rest-api-template/data/presenter"
	"github.com/jesseokeya/go-rest-api-template/server/api"
	"github.com/upper/db/v4"
)

type signupRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Invite    string `json:"invite"`
}

func (u *signupRequest) Bind(r *http.Request) error {
	if u.FirstName == "" || u.LastName == "" {
		return errors.New("missing name")
	}
	return emailx.Validate(u.Email)
}

// Signup creates a new user
func Signup(w http.ResponseWriter, r *http.Request) {
	newSignup := &signupRequest{}
	if err := render.Bind(r, newSignup); err != nil {
		api.IgnoreError(render.Render(w, r, api.ErrInvalidEmailSignup(err)))
		return
	}

	// encrypt with bcrypt
	epw, err := encrypt(newSignup.Password)
	if err != nil {
		// mask the encryption error and return
		api.IgnoreError(render.Render(w, r, api.ErrEncryptionError))
		return
	}

	// checking if the user with this email already exists.
	// the source does NOT have to be email.
	//  -> this prevents duplicated log-in accounts with similar email.
	user, err := data.DB.User.FindOne(db.Cond{
		"email": emailx.Normalize(newSignup.Email),
	})
	if err != nil && err != db.ErrNoMoreRows {
		api.IgnoreError(render.Render(w, r, api.ErrDatabase(err)))
		return
	}

	// throws error if user already exists
	if user != nil {
		api.IgnoreError(render.Render(w, r, api.ErrUserExists))
		return
	}

	// necesarry fields required to create a new user
	user = &data.User{
		FirstName:    newSignup.FirstName,
		LastName:     newSignup.LastName,
		Email:        emailx.Normalize(newSignup.Email),
		PasswordHash: string(epw),
	}

	if err := data.DB.Save(user); err != nil {
		api.IgnoreError(render.Render(w, r, api.ErrDatabase(err)))
		return
	}

	presented := presenter.NewAuthUser(r.Context(), user)
	api.Render(w, r, presented)
}
