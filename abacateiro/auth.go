package application

import (
	"context"
	"fmt"
	"unicode/utf8"
)

const (
	MinUsernameLen = 4
	MaxUsernameLen = 72
	MinPasswordLen = 8
	MaxPasswordLen = 72
)

type LoginUserQuery struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthInfo struct {
	AuthId string
	Email  string
}

// Validate returns an error if LoginUserQUery has invalid fields.
// Only performs basic validation.
func (l *LoginUserQuery) Validate() error {
	if l.Username == "" {
		return fmt.Errorf("username must be provided")
	} else if utf8.RuneCountInString(l.Username) > MaxUsernameLen {
		return fmt.Errorf("username too long")
	} else if utf8.RuneCountInString(l.Username) < MinUsernameLen {
		return fmt.Errorf("username too short")
	} else if l.Password == "" {
		return fmt.Errorf("password must be provided")
	} else if utf8.RuneCountInString(l.Password) > MaxPasswordLen {
		return fmt.Errorf("password too long")
	} else if utf8.RuneCountInString(l.Password) < MinPasswordLen {
		return fmt.Errorf("password too short")
	}
	return nil
}

type AuthService interface {
	Login(ctx context.Context, query *LoginUserQuery) (*AuthInfo, error)
	// Logout(ctx context.Context, authId string) error
}
