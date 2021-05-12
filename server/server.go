package server

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/jesseokeya/go-rest-api-template/data"
	"github.com/jesseokeya/go-rest-api-template/lib/session"
	"github.com/jesseokeya/go-rest-api-template/server/api"
	"github.com/jesseokeya/go-rest-api-template/server/auth"
	"github.com/jesseokeya/go-rest-api-template/server/oauth"
	"github.com/jesseokeya/go-rest-api-template/server/user"
	"github.com/rs/zerolog"
)

// Server holds refernces to other interfaces to be used
type Server struct {
	lg *zerolog.Logger
	o  *oauth.OAuth

	opts ServerOptions
}

type ServerOptions struct {
	debug bool
}

type ServerOption func(*ServerOptions) error

func Debug(b bool) ServerOption {
	return func(o *ServerOptions) error {
		o.debug = b
		return nil
	}
}

// New initializes the server
func New(auth *session.Auth, db *data.Database, options ...ServerOption) (*Server, error) {
	logger := zerolog.New(os.Stdout)
	logger = logger.Hook(api.SeverityHook{
		AlertFn: func(level zerolog.Level, msg string) {},
	})
	s := &Server{
		lg: &logger,
	}

	for _, opt := range options {
		if err := opt(&s.opts); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	api.IgnoreError(w.Write([]byte(`👾`)))
}

func (s *Server) NewStructuredLogger() func(next http.Handler) http.Handler {
	if s.opts.debug {
		return middleware.Logger
	}
	return middleware.RequestLogger(&api.StructuredLogger{Logger: s.lg})
}

// Routes initializes the middlewares and routes to be used in the application
func (s *Server) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.NoCache)
	r.Use(middleware.RequestID)
	r.Use(s.NewStructuredLogger())
	r.Use(middleware.Recoverer)
	r.Use(api.Cors().Handler)

	r.Get("/", healthCheck)

	r.Group(func(r chi.Router) {
		r.Mount("/auth", s.o.Routes())
	})

	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(s.o.Authority()))

		// Handle valid / invalid tokens.
		r.Use(auth.Authenticator)

		// User Session context
		r.Use(auth.SessionCtx)

		r.Mount("/users", user.Routes())
	})

	return r
}
