package oauth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/jesseokeya/go-rest-api-template/data"
	"github.com/jesseokeya/go-rest-api-template/data/presenter"
	"github.com/jesseokeya/go-rest-api-template/lib/session"
	"github.com/jesseokeya/go-rest-api-template/server/api"
)

type OAuth struct {
	*session.Auth

	srv *server.Server
}

func New(tokAuth *session.Auth, db *data.Database) *OAuth {
	// [oauth manager]
	manager := manage.NewDefaultManager()

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	return &OAuth{
		Auth: tokAuth,
		srv:  srv,
	}
}

var (
	ErrUserID = errors.New("invalid or missing user id")
)

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		fmt.Fprintf(w, "Permission: %v", err)
		return
	}
	userId, ok := token.Get("userId")
	if !ok {
		return "", ErrUserID
	}
	return fmt.Sprintf("%v", userId), nil
}

func (o *OAuth) Authorize(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		api.IgnoreError(render.Render(w, r, api.ErrInvalidRequest(err)))
		return
	}
	err := o.srv.HandleAuthorizeRequest(w, r)
	if err != nil {
		api.IgnoreError(render.Render(w, r, api.ErrInvalidRequest(err)))
		return
	}
}

func (o *OAuth) GetMe(w http.ResponseWriter, r *http.Request) {
	token, err := o.srv.ValidationBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := strconv.ParseInt(token.GetUserID(), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := data.DB.User.FindByID(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.IgnoreError(json.NewEncoder(w).Encode(presenter.NewUser(r.Context(), user)))
}

func (o *OAuth) GetToken(w http.ResponseWriter, r *http.Request) {
	api.IgnoreError(o.srv.HandleTokenRequest(w, r))
}
