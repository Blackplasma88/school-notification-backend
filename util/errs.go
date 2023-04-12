package util

import "errors"

var (
	ErrNotFound                 = errors.New("not found")
	ErrEventInvalid             = errors.New("event invalid")
	ErrIdIsNotPrimitiveObjectID = errors.New("Id is not primitive objectID")
	ErrRequireParameter         = errors.New("require parameter ")
	ErrInternalServerError      = errors.New("internal server error")
	ErrValueAlreadyExists       = errors.New(" does already exists")
	ErrValueNotAlreadyExists    = errors.New(" does not already exists")
	ErrValueInvalid             = errors.New(" value is invalid")
	// ErrValueNotMatch             = errors.New(" value not match")
	ErrTypeInvalid               = errors.New("type is invalid")
	ErrProfileIdAlreadyExists    = errors.New("The profile id already exists")
	ErrProfileIdNotAlreadyExists = errors.New("The profile id does not already exists")
	ErrStatusInvalid             = errors.New(" status invalid expect: ")
)

func ReturnError(err string) error {
	return errors.New(err)
}

func ReturnErrorStatusInvalid(name string, expect string) error {
	return errors.New(name + ErrStatusInvalid.Error() + expect)
}
