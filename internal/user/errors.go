package user

import "errors"

var (
	ErrNotFound      = errors.New("Benutzer nicht gefunden")
	ErrAlreadyExists = errors.New("Benutzer existiert bereits")
)
