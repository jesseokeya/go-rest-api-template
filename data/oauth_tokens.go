package data

import (
	"context"
	"errors"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

type OAuthToken struct {
	ID        int64          `db:"id,omitempty"`
	CreatedAt *time.Time     `db:"created_at"`
	ExpiresAt time.Time      `db:"expires_at"`
	Code      string         `db:"code"`
	Access    string         `db:"access"`
	Refresh   string         `db:"refresh"`
	Data      OAuthTokenData `db:"data"`
}

type OAuthTokenData struct {
	*models.Token
	*postgresql.JSONBConverter
}

var _ = interface {
	db.Record
	db.BeforeCreateHook
}(&OAuthToken{})

var (
	ErrInvalidOAuthInfo = errors.New("invalid oauth TokenInfo type")
)

var (
	maxActiveToken = 10
)

// BeforeCreate retrives and sets GetTimeUTCPointer before creating a new oauth token
func (u *OAuthToken) BeforeCreate(sess db.Session) error {
	u.CreatedAt = GetTimeUTCPointer()
	return nil
}

func (u *OAuthToken) Store(sess db.Session) db.Store {
	return OAuthTokens(sess)
}

// OAuthTokensStore holds oauth_tokens database collection
type OAuthTokensStore struct {
	db.Collection
}

var _ = interface {
	db.Store
}(&OAuthTokensStore{})

func OAuthTokens(sess db.Session) *OAuthTokensStore {
	return &OAuthTokensStore{sess.Collection("oauth_tokens")}
}

// create and store the new token information
func (s *OAuthTokensStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	tok, ok := info.(*models.Token)
	if !ok {
		return ErrInvalidOAuthInfo
	}

	item := &OAuthToken{Data: OAuthTokenData{Token: tok}}
	if code := info.GetCode(); code != "" {
		item.Code = code
		item.ExpiresAt = info.GetCodeCreateAt().Add(info.GetCodeExpiresIn())
	} else {
		item.Access = info.GetAccess()
		item.ExpiresAt = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn())
		if refresh := info.GetRefresh(); refresh != "" {
			item.Refresh = info.GetRefresh()
			item.ExpiresAt = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn())
		}

		// Limit number of similar and active tokens to 10.
		findCond := db.Cond{db.Raw("data->>'UserID'"): info.GetUserID()}
		findQ := s.Find(findCond).OrderBy("expires_at")
		var tokens []*OAuthToken
		if err := findQ.All(&tokens); err != nil {
			return err
		}
		if len(tokens) == maxActiveToken {
			// max limit reached, remove the oldest
			if err := DB.Delete(tokens[0]); err != nil {
				return err
			}
		}
	}

	return DB.Save(item)
}

// delete the authorization code
func (s *OAuthTokensStore) RemoveByCode(ctx context.Context, code string) error {
	return s.Find(db.Cond{"code": code}).Delete()
}

// use the access token to delete the token information
func (s *OAuthTokensStore) RemoveByAccess(ctx context.Context, access string) error {
	return s.Find(db.Cond{"access": access}).Delete()
}

// use the refresh token to delete the token information
func (s *OAuthTokensStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return s.Find(db.Cond{"refresh": refresh}).Delete()
}

// use the authorization code for token information data
func (s *OAuthTokensStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	return s.FindOne(db.Cond{"code": code})
}

// use the access token for token information data
func (s *OAuthTokensStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	return s.FindOne(db.Cond{"access": access})
}

// use the refresh token for token information data
func (s *OAuthTokensStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	return s.FindOne(db.Cond{"refresh": refresh})
}

// find one token
func (s *OAuthTokensStore) FindOne(cond db.Cond) (oauth2.TokenInfo, error) {
	var token OAuthToken
	if err := s.Find(cond).One(&token); err != nil {
		return nil, err
	}
	return token.Data, nil
}

// find all token
func (s *OAuthTokensStore) FindAll(cond db.Cond) ([]*OAuthToken, error) {
	var tokens []*OAuthToken
	if err := s.Find(cond).All(&tokens); err != nil {
		return nil, err
	}
	return tokens, nil
}
