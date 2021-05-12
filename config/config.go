package config

import (
	"fmt"
	"os"

	"github.com/burntsushi/toml"
	"github.com/jesseokeya/go-rest-api-template/data"
	"github.com/jesseokeya/go-rest-api-template/lib/connect"
	"github.com/jesseokeya/go-rest-api-template/lib/session"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// Config holds all the configuration fields needed within the application
type Config struct {
	Environment string `toml:"environment"`
	Bind        string `toml:"bind"`

	// [db]
	DB data.DBConf `toml:"db"`

	// [connect]
	Connect connect.Configs `toml:"connect"`

	// [jwt]
	JWT session.Config `toml:"jwt"`

	// Heroku requires env PORT to set its internal port for running our application
	Port string
}

// NewFromFile instantiates the config struct
func NewFromFile(fileConfig, envConfig string) (*Config, error) {
	file := fileConfig
	if file == "" {
		file = envConfig
	}

	var conf Config
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, err
	}
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		return nil, fmt.Errorf("unable to load config file: %w", err)
	}

	// If PORT environment variable is set, use it.
	if v := os.Getenv("PORT"); v != "" {
		conf.Port = v
	}

	// If development, set zerolog to debug
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	if conf.Environment != "production" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if conf.Environment == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return &conf, nil
}
