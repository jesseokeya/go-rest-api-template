package api

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type ContextKey struct {
	Name string
}

func (k *ContextKey) String() string {
	return "chi context value " + k.Name
}

var (
	// SessionCtxKey is the context.Context key to store the request context.
	SessionUserCtxKey = &ContextKey{Name: "Session.User"}
	SessionRoleCtxKey = &ContextKey{Name: "Session.Role"}
)
