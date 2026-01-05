package auth

import (
	"errors"
)

// custom errors for better error handling
var (
	ErrPasswordTooShort = errors.New("password must be at least 8 characters long")
	ErrPasswordTooWeek  = errors.New("password is too week")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrHashingPassword  = errors.New("failed to hash password")
)

type PasswordStrength int

const (
	StrengthAny PasswordStrength = iota
	StrengthBasic
	StrengthStrong
)

type PasswordValidator struct {
	minLength      int
	requireUpper   bool
	requireLower   bool
	requireNumber  bool
	requireSpecial bool
}

func NewPasswordValidator(minLength int, strength PasswordStrength) *PasswordValidator {
	pv := &PasswordValidator{
		minLength: minLength,
	}

	switch strength {
	case StrengthBasic:
		pv.requireLower = true
		pv.requireUpper = true
		pv.requireNumber = true
	case StrengthStrong:
		pv.requireLower = true
		pv.requireUpper = true
		pv.requireSpecial = true
		pv.requireNumber = true
	}

	return pv
}
