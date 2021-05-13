package data

import (
	"context"
	"errors"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

type OAuthClient struct {
	ID     string          `db:"id,omitempty"`
	Secret string          `db:"secret"`
	Domain string          `db:"domain"`
	Data   OAuthClientData `db:"data"`
}

type OAuthClientData struct {
	*models.Client
	*postgresql.JSONBConverter
}

var (
	ErrInvalidOAuthClient = errors.New("invalid oauth client type")
)

func (u *OAuthClient) Store(sess db.Session) db.Store {
	return OAuthClients(sess)
}

// OAuthClientsStore holds oauth_tokens database collection
type OAuthClientsStore struct {
	db.Collection
}

func OAuthClients(sess db.Session) *OAuthClientsStore {
	return &OAuthClientsStore{sess.Collection("oauth_clients")}
}

// according to the ID for the client information
func (s *OAuthClientsStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	return s.FindOne(db.Cond{"id": id})
}

// find one client
func (s *OAuthClientsStore) FindOne(cond db.Cond) (oauth2.ClientInfo, error) {
	var client OAuthClient
	if err := s.Find(cond).One(&client); err != nil {
		return nil, err
	}
	return client.Data, nil
}

// Create creates and stores the new client information
func (s *OAuthClientsStore) Create(info oauth2.ClientInfo) error {
	client, ok := info.(*models.Client)
	if !ok {
		return ErrInvalidOAuthClient
	}
	return DB.Save(&OAuthClient{
		ID:     info.GetID(),
		Secret: info.GetSecret(),
		Domain: info.GetDomain(),
		Data:   OAuthClientData{Client: client},
	})
}
