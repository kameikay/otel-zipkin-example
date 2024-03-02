package exceptions

import "errors"

var (
	ErrInvalidCEP        = errors.New("invalid zipcode")
	ErrCannotFindZipcode = errors.New("can not find zipcode")
)
