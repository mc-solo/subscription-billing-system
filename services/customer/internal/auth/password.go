package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
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
		return fmt.Errorf("%w: minimun %d characters required", ErrPasswordTooShort, pv.minLength)
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
			return fmt.Errorf("%w: must contain at least one number", ErrPasswordTooWeak)
		}
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

func HashPassword(password string, cost int) (string, error) {
	// validate cost param [4-31]
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return "", fmt.Errorf("bcrypt cost must be between %d and %d", bcrypt.MinCost, bcrypt.MaxCost)
	}

	passwordBytes := []byte(password)

	hashedBytes, err := bcrypt.GenerateFromPassword(passwordBytes, cost)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrHashingPassword, err)
	}

	// convert bytes back to strings for storage
	return string(hashedBytes), nil
}

// compares a plaintext password with bcrypt hash
func CheckPassword(password, hashedPassword string) error {
	// convert bytes back to arrays
	passwordBytes := []byte(password)
	hashedBytes := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(hashedBytes, passwordBytes)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidPassword, err)
	}
	return nil
}

// generate random pass [useful for initial user setup]
func GenerateRandomPassword(length int) (string, error) {
	if length < 8 {
		return "", errors.New("password length must be at least 8 characters long")
	}

	// char sets
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lower := "abcdefghijklmnopqrstuvwxyz"
	digits := "0123456789"
	special := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	all := upper + lower + digits + special

	bytes := make([]byte, length)

	// read rand bytes
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random password: %w", err)
	}

	// covnert bytes to chars
	password := make([]byte, length)
	for i := range password {
		password[i] = all[int(bytes[i])%len(all)]
	}

	// ensures at least one of each char type
	password[0] = upper[int(bytes[0])%len(upper)]
	password[1] = lower[int(bytes[1])%len(lower)]
	password[2] = digits[int(bytes[2])%len(digits)]
	password[3] = special[int(bytes[3])%len(special)]

	// shuffle the password
	rand.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password), nil
}

// extracts info from bcrypt hash
func GenerateHashInfo(hashedPassword string) (cost int, err error) {
	if len(hashedPassword) < 7 || !strings.HasPrefix(hashedPassword, "$2a$") {
		return 0, errors.New("invalid bcrypt hash format")
	}

	// extract the cost part (bytes 4-6: "12$")
	costStr := hashedPassword[4:6]
	cost, err = strconv.Atoi(costStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse bcrypt cost: %w", err)
	}

	return cost, nil
}
