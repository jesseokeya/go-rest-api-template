package data

import (
	"time"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

// User holds postgres data structure for user
type User struct {
	ID           int64      `db:"id,omitempty" json:"id"`
	Email        string     `db:"email,omitempty" json:"email"`
	PasswordHash string     `db:"password_hash,omitempty" json:"-"`
	Role         UserRole   `db:"role" json:"role"`
	FirstName    string     `db:"first_name,omitempty" json:"firstName"`
	LastName     string     `db:"last_name,omitempty" json:"lastName"`
	Location     string     `db:"location,omitempty" json:"location"`
	Etc          UserEtc    `db:"etc,omitempty" json:"etc"`
	CreatedAt    *time.Time `db:"created_at,omitempty" json:"createdAt"`
	UpdatedAt    *time.Time `db:"updated_at,omitempty" json:"updatedAt"`

	// Display Only
	Biography string `json:"biography"`
	Name      string `json:"name"`
}

// UserRole is users role
type UserRole string

const (
	// UserRoleMember is user role member
	UserRoleMember UserRole = "member"
	// UserRoleAdmin is user role admin
	UserRoleAdmin UserRole = "admin"
	// UserRoleBlocked is user role blocked
	UserRoleBlocked UserRole = "blocked"
)

// UserEtc is a collection of auxillary key/values around a user
type UserEtc struct {
	HasProfile bool

	*postgresql.JSONBConverter
}

var _ = interface {
	db.Record
	db.BeforeUpdateHook
	db.BeforeCreateHook
}(&User{})

// Store initializes a session interface that defines methods for database adapters.
func (u *User) Store(sess db.Session) db.Store {
	return Users(sess)
}

// BeforeCreate retrives and sets GetTimeUTCPointer before creating a new user
func (u *User) BeforeCreate(sess db.Session) error {
	if err := u.BeforeUpdate(sess); err != nil {
		return err
	}
	u.UpdatedAt = nil
	u.CreatedAt = GetTimeUTCPointer()
	return nil
}

// BeforeUpdate retrives and sets GetTimeUTCPointer before updating a user
func (u *User) BeforeUpdate(sess db.Session) error {
	u.UpdatedAt = GetTimeUTCPointer()
	return nil
}

// UsersStore holds thread users database collection
type UsersStore struct {
	db.Collection
}

var _ = interface {
	db.Store
}(&UsersStore{})

// Users retrieves a list of all users
func Users(sess db.Session) *UsersStore {
	return &UsersStore{sess.Collection("users")}
}

// FindByID retrieves a particular user by id
func (store UsersStore) FindByID(ID int64) (*User, error) {
	return store.FindOne(db.Cond{"id": ID})
}

// FindByAlias retrieves a particular user by alias
func (store UsersStore) FindByAlias(alias string) (*User, error) {
	return store.FindOne(db.Cond{"alias": alias})
}

// FindByEmail retrieves a particular user by email
func (store UsersStore) FindByEmail(email string) (*User, error) {
	return store.FindOne(db.Cond{"email": email})
}

// FindOne retrieves user by certain conditions from UsersStore
func (store UsersStore) FindOne(cond ...interface{}) (*User, error) {
	var user *User
	if err := store.Find(cond...).One(&user); err != nil {
		return nil, err
	}
	return user, nil
}

// FindAll retrieves users by certain conditions from UsersStore
func (store UsersStore) FindAll(cond ...interface{}) ([]*User, error) {
	var users []*User
	if err := store.Find(cond...).All(&users); err != nil {
		return nil, err
	}
	return users, nil
}
