package ucli

import (
	"errors"
	"strings"

	"github.com/charmbracelet/huh"
)

type CreateUserInput struct {
	UID         string
	DisplayName string
	FirstName   string
	LastName    string
}

func promptCreateUser() (*CreateUserInput, error) {
	var in CreateUserInput

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("UID (login secret)").
				Value(&in.UID).
				Validate(nonEmpty("uid")),

			huh.NewInput().
				Title("Display name").
				Value(&in.DisplayName).
				Validate(nonEmpty("display name")),

			// optional by default
			huh.NewInput().
				Title("First name").
				Value(&in.FirstName),

			// optional by default
			huh.NewInput().
				Title("Last name").
				Value(&in.LastName),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}

	return &in, nil
}

func nonEmpty(label string) func(string) error {
	return func(v string) error {
		if strings.TrimSpace(v) == "" {
			return errors.New(label + " is required")
		}
		return nil
	}
}
