package validation

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile("^[a-z0-9_]+$")
	isValidFullname = regexp.MustCompile("^[a-zA-Z\\s]+$")
)

func ValidateString(val string, minLength, maxLength int) error {
	n := len(val)
	if n > maxLength || n < minLength {
		return fmt.Errorf("must contain %d-%d characters", minLength, maxLength)
	}

	return nil
}

func ValidateUsername(val string) error {
	if err := ValidateString(val, 3, 60); err != nil {
		return err
	}

	if !isValidUsername.MatchString(val) {
		return fmt.Errorf("must contains only lowercase letters, digits, or underscore")
	}

	return nil
}

func ValidatePassword(val string) error {
	return ValidateString(val, 6, 100)
}

func ValidateEmail(val string) error {
	if err := ValidateString(val, 5, 200); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(val); err != nil {
		return fmt.Errorf("must contain valid email address")
	}

	return nil
}

func ValidateFullname(val string) error {
	if err := ValidateString(val, 3, 100); err != nil {
		return err
	}

	if !isValidFullname.MatchString(val) {
		return fmt.Errorf("must contains only letters, or spaces")
	}

	return nil
}
