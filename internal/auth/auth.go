package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

func (a *Auth) Login(ctx *gin.Context, uid string) (*user.User, error) {
	u, err := a.users.FindByUID(ctx.Request.Context(), uid)
	if err != nil {
		return nil, err // z. B. ErrNotFound
	}

	sess := sessions.Default(ctx)
	sess.Set(SessionKeyUID, u.UID)

	if err := sess.Save(); err != nil {
		return nil, err
	}

	return u, nil
}

func (a *Auth) Logout(ctx *gin.Context) error {
	sess := sessions.Default(ctx)

	// sess.Delete(SessionKeyUID)
	sess.Set(SessionKeyUID, "")

	return sess.Save()
}

func (a *Auth) IsAuthenticated(ctx *gin.Context) bool {
	_, ok := a.CurrentUser(ctx)
	return ok
}

func (a *Auth) CurrentUser(ctx *gin.Context) (*user.User, bool) {
	sess := sessions.Default(ctx)

	rawUID := sess.Get(SessionKeyUID)
	if rawUID == nil {
		return nil, false
	}

	uid, ok := rawUID.(string)
	if !ok || uid == "" {
		return nil, false
	}

	u, err := a.users.FindByUID(ctx.Request.Context(), uid)
	if err != nil {
		return nil, false
	}

	return u, true
}
