package schema

import (
	"errors"
	"regexp"
)

var _ V = new(IsPhoneNumber)
var _ RealV = new(phoneNumber)

type IsPhoneNumber struct {

}

func (v *IsPhoneNumber) INJ() RealV {
	return &phoneNumber{}
}

type phoneNumber struct {
	val string
}

var c = regexp.MustCompile(`^1\d{10}$`)

func (pn *phoneNumber) Check() error {
	if pn.val != "" && c.MatchString(pn.val){
		return nil
	}
	return errors.New("bad")
}
