package util

import "errors"

var (
	ErrNotFound                 = errors.New("not found")
	ErrIdIsNotPrimitiveObjectID = errors.New("Id is not primitive objectID")
	ErrRequireParameter         = errors.New("require parameter ")
	ErrInternalServerError      = errors.New("internal server error")
	ErrValueAlreadyExists       = errors.New(" does already exists")
	ErrValueNotAlreadyExists    = errors.New(" does not already exists")
	ErrValueInvalid             = errors.New(" value is invalid")
)

func ReturnError(err string) error {
	return errors.New(err)
}
