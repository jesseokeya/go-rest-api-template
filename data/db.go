package data

import (
	"fmt"
	"strings"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

// Database holds database postgres table structure
type Database struct {
	db.Session

	User *UsersStore
}

// DBConf database configuration
type DBConf struct {
	Database        string   `toml:"database"`
	Hosts           []string `toml:"hosts"`
	Username        string   `toml:"username"`
	Password        string   `toml:"password"`
	DebugQueries    bool     `toml:"debug_queries"`
	ApplicationName string   `toml:"application_name"`
	MaxConnection   int      `toml:"max_connection"`
	SSLMode         string   `toml:"ssl_mode"`
	DatabaseURL     string   `env:"DATABASE_URL"`
}

// DB is a pointer to the database interface
var DB *Database

// String converts database config to a connection url
func (cf *DBConf) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cf.Username, cf.Password, strings.Join(cf.Hosts, ","), cf.Database, cf.SSLMode)
}

// NewDB initializes / opens a new database connection
func NewDB(conf DBConf) (*Database, error) {
	connString := conf.String()

	if conf.DatabaseURL != "" {
		connString = conf.DatabaseURL
	}

	connURL, err := postgresql.ParseURL(connString)
	if err != nil {
		return nil, err
	}
	// extra options
	connURL.Options["application_name"] = conf.ApplicationName

	db := &Database{}
	db.Session, err = postgresql.Open(connURL)
	if err != nil {
		return nil, err
	}

	db.User = Users(db.Session)

	// global instance for access across modules
	DB = db

	return db, nil
}
