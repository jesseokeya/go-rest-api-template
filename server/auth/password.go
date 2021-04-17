package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/jesseokeya/go-rest-api-template/data"
	"github.com/jesseokeya/go-rest-api-template/lib/connect"
	"github.com/jesseokeya/go-rest-api-template/server/api"
	"github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordLength int = 8
	bCryptCost        int = 10

	// if the user's password hash is empty, use this
	// hash to mask the fact
	timingHash = "$2a$10$4Kys.PIxpCIoUmlcY6D7QOTuMPgk27lpmV74OWCWfqjwnG/JN4kcu"
)

var (
	ErrPasswordResetExpired = errors.New("password reset link is expired")
)

// bcrypt compare hash with given password
func VerifyPassword(hash, password string) bool {
	// incase either hash or password is empty, compare
	// something and return false to mask the timing
	if len(hash) == 0 || len(password) == 0 {
		if err := bcrypt.CompareHashAndPassword([]byte(timingHash), []byte(password)); err != nil {
			return false
		}
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func encrypt(password string) ([]byte, error) {
	// encrypt with bcrypt
	return bcrypt.GenerateFromPassword([]byte(password), bCryptCost)
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

func (p *forgotPasswordRequest) Bind(r *http.Request) error {
	if p.Email == "" {
		return errors.New("empty email")
	}
	return nil
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var payload forgotPasswordRequest
	if err := render.Bind(r, &payload); err != nil {
		api.IgnoreError(render.Render(w, r, api.ErrInvalidRequest(err)))
		return
	}

	user, err := data.DB.User.FindByEmail(payload.Email)
	if err != nil {
		api.IgnoreError(render.Render(w, r, api.ErrInvalidRequest(err)))
		return
	}

	if user.Etc.OnetimeCodeExpire > time.Now().Unix() {
		// there is a code that's still active. ignore
		w.WriteHeader(http.StatusNotModified)
		return
	}

	t := fmt.Sprintf("%d", time.Now().UnixNano())
	mac := hmac.New(sha256.New, []byte(t))
	mac.Write([]byte(data.RandString(7)))
	user.Etc.OnetimeCode = hex.EncodeToString(mac.Sum(nil))
	user.Etc.OnetimeCodeExpire = time.Now().Add(24 * time.Hour).Unix()

	if err := data.DB.Save(user); err != nil {
		api.Render(w, r, api.ErrServiceUnavailable(err))
		return
	}

	ctx := r.Context()
	if err := connect.SD.SendPasswordReset(ctx, user); err != nil {
		api.Render(w, r, api.ErrServiceUnavailable(err))
		return
	}

	if err := connect.SG.OnetimeCode(user); err != nil {
		api.Render(w, r, api.ErrServiceUnavailable(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.Respond(w, r, "ok")
}

type resetPasswordRequest struct {
	OnetimeCode     string `json:"onetimeCode"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func (p *resetPasswordRequest) Bind(r *http.Request) error {
	if p.OnetimeCode == "" {
		return errors.New("invalid request")
	}
	if p.Password == "" || p.ConfirmPassword == "" {
		return errors.New("invalid password")
	}
	if p.Password != p.ConfirmPassword {
		return errors.New("mismatch password")
	}
	return nil
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var payload resetPasswordRequest
	if err := render.Bind(r, &payload); err != nil {
		api.IgnoreError(render.Render(w, r, api.ErrInvalidRequest(err)))
		return
	}

	user, err := data.DB.User.FindOne(db.Cond{db.Raw("etc->>'otc'"): payload.OnetimeCode})
	if err != nil {
		api.IgnoreError(render.Render(w, r, api.ErrInvalidRequest(err)))
		return
	}

	// check if onetimecode is expired.
	if user.Etc.OnetimeCodeExpire < time.Now().Unix() {
		api.IgnoreError(render.Render(w, r, api.ErrInvalidRequest(ErrPasswordResetExpired)))
		return
	}

	user.Etc.OnetimeCode = ""
	user.Etc.OnetimeCodeExpire = 0
	epw, err := encrypt(payload.Password)
	if err != nil {
		// mask the encryption error and return
		api.IgnoreError(render.Render(w, r, api.ErrEncryptionError))
		return
	}
	user.PasswordHash = string(epw)

	if err := data.DB.Save(user); err != nil {
		api.IgnoreError(render.Render(w, r, api.ErrServiceUnavailable(err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.Respond(w, r, "ok")
}
