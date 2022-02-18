package domainerror

import "errors"

var (
	ErrNotExists = errors.New("not exists")
	ErrNilObject = errors.New("nil object")
)
