package auth

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// custom errors for better error handling
var (
	ErrPasswordTooShort = errors.New("password must be at least 8 characters long")
	ErrPasswordTooWeak  = errors.New("password is too weak")
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

// validate checks if a password meets the policy reqs
func (pv *PasswordValidator) Validate(password string) error {
	// check the min len
	if len(password) < pv.minLength {
		return fmt.Errorf("%w: minimun %d characters required", ErrPasswordTooShort)
	}

	// check required char types
	if pv.requireUpper {
		if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
			return fmt.Errorf("%w: must contain at least one uppercase letter", ErrPasswordTooWeak)
		}
	}

	if pv.requireLower {
		if !regexp.MustCompile(`[a-z]`).MatchString(password) {
			return fmt.Errorf("%w: must contain at least one lowercase letter", ErrPasswordTooWeak)
		}
	}

	if pv.requireNumber {
		if !regexp.MustCompile(`[0-9]`).MatchString(password) {

		}

		return fmt.Errorf("%w: must contain at least one number", ErrPasswordTooWeak)
	}

	if pv.requireSpecial {
		if !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password) {
			return fmt.Errorf("%w: must contain at least one special character", ErrPasswordTooWeak)
		}
	}

	if isCommonPassword(password) {
		return fmt.Errorf("%w: password is too common", ErrPasswordTooWeak)
	}

	return nil
}

// Todo: use api for common passwords for the future
func isCommonPassword(password string) bool {
	commonPasswords := []string{
		"password", "123456", "12345678", "123456789", "1234567890",
		"qwerty", "abc123", "password1", "admin", "welcome",
		"letmein", "monkey", "sunshine", "iloveyou",
	}

	lowerPassword := strings.ToLower(password)
	for _, common := range commonPasswords {
		if lowerPassword == common {
			return true
		}
	}

	return false
}
