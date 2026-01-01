package auth

import (
	"net/http"

	"github.com/axelrhd/hagg/internal/session"
	"github.com/axelrhd/hagg/internal/user"
)

const SessionKeyUID = "uid"

type Auth struct {
	users user.Store
}

func New(users user.Store) *Auth {
	return &Auth{
		users: users,
	}
}

// Login authenticates a user and creates a session.
func (a *Auth) Login(req *http.Request, uid string) (*user.User, error) {
	u, err := a.users.FindByUID(req.Context(), uid)
	if err != nil {
		return nil, err
	}

	session.Manager.Put(req.Context(), SessionKeyUID, u.UID)
	return u, nil
}

// Logout clears the user session.
func (a *Auth) Logout(req *http.Request) error {
	session.Manager.Put(req.Context(), SessionKeyUID, "")
	return nil
}

// IsAuthenticated checks if a user is authenticated.
func (a *Auth) IsAuthenticated(req *http.Request) bool {
	_, ok := a.CurrentUser(req)
	return ok
}

// CurrentUser retrieves the currently authenticated user.
func (a *Auth) CurrentUser(req *http.Request) (*user.User, bool) {
	rawUID := session.Manager.Get(req.Context(), SessionKeyUID)
	uid, ok := rawUID.(string)
	if !ok || uid == "" {
		return nil, false
	}

	u, err := a.users.FindByUID(req.Context(), uid)
	if err != nil {
		return nil, false
	}

	return u, true
}
