package schema

import (
	"errors"
	"regexp"
)

var c = regexp.MustCompile(`^1\d{10}$`)

type IsPhoneNumber struct{}

func (s IsPhoneNumber) check(value string, metadata string) error {
	if value != "" && c.MatchString(value) {
		return nil
	}
	return errors.New("isnt phone number")
}
