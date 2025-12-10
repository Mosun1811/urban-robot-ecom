package service

import (
	"errors"
	"regexp"
)

// -------------------------
// NAME VALIDATION
// -------------------------
var nameRegex = regexp.MustCompile(`^[A-Za-z\s\-]{2,50}$`)

func ValidateName(name string) error {
	if !nameRegex.MatchString(name) {
		return errors.New("name must be 2–50 characters and only letters, spaces, or hyphens")
	}
	return nil
}

// -------------------------
// EMAIL VALIDATION
// -------------------------
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// -------------------------
// PASSWORD VALIDATION
// -------------------------
func ValidatePassword(pw string) error {
	if len(pw) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(pw) > 64 {
		return errors.New("password cannot exceed 64 characters")
	}

	uppercase := regexp.MustCompile(`[A-Z]`)
	lowercase := regexp.MustCompile(`[a-z]`)
	number := regexp.MustCompile(`[0-9]`)
	special := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]{};':"\\|,.<>\/?]`)

	if !uppercase.MatchString(pw) {
		return errors.New("password must contain at least 1 uppercase letter")
	}
	if !lowercase.MatchString(pw) {
		return errors.New("password must contain at least 1 lowercase letter")
	}
	if !number.MatchString(pw) {
		return errors.New("password must contain at least 1 number")
	}
	if !special.MatchString(pw) {
		return errors.New("password must contain at least 1 special character")
	}

	return nil
}
