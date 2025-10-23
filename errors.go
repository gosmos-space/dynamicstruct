package dynamicstruct

import "errors"

var (
	ErrFieldAlreadyExists   = errors.New("field already exists")
	ErrInstanceAlreadyBuilt = errors.New("instance already built")
	ErrInstanceNotBuilt     = errors.New("instance not built")
	ErrValueMustBePointer   = errors.New("value must be a pointer")
	ErrValueCannotBeNil     = errors.New("value cannot be nil")
	ErrFieldNotFound        = errors.New("field not found")
	ErrIncompatibleTypes    = errors.New("incompatible types of value and field")
	ErrInvalidTag           = errors.New("invalid struct tag format")
)
