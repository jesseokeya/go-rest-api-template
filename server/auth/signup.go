package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/clout-jam/api/data"
	"github.com/clout-jam/api/data/presenter"
	"github.com/clout-jam/api/lib/connect"
	"github.com/clout-jam/api/server/api"
	"github.com/go-chi/render"
	"github.com/goware/emailx"
	"github.com/rs/zerolog/log"
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
		Status:       data.UserStatusIncomplete,
		PasswordHash: string(epw),
	}

	// Check if this user is a invite
	if code := newSignup.Invite; code != "" {
		func(code string, u *data.User) {
			invite, err := data.DB.Invite.FindByCode(code)
			if err != nil {
				log.Error().Msgf("User signup invite with code %s: %v", code, err)
				return
			}

			// Track who the inviter was
			u.InviterID = &(invite.InviterID)
			// Track when this "accepted" happened
			invite.AcceptedAt = data.GetTimeUTCPointer()
			_ = data.DB.Save(invite)

			// Send a slack notification under the same thread as the original inviter's review request
			inviter, _ := data.DB.User.FindByID(invite.InviterID)
			connect.SL.SendThreadReply("signup", fmt.Sprintf("%s invite just signed up!", u.Email), inviter.Etc.ReviewID)
		}(code, user)
	}

	if err := data.DB.Save(user); err != nil {
		api.IgnoreError(render.Render(w, r, api.ErrDatabase(err)))
		return
	}

	presented := presenter.NewAuthUser(r.Context(), user)
	presented.Intercom = &presenter.AuthIntercom{
		HashID: connect.SG.GetUserKey(user.ID),
		Hash:   connect.SG.GetUserHash(user.ID),
	}
	api.Render(w, r, presented)
}
