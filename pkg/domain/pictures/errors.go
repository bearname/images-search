package pictures

import "errors"

var (
	ErrNotFound         = errors.New("pictures not exist")
	ErrCannotAssignNull = errors.New("cannot assign NULL to *string")
	ErrEmptyImages      = errors.New("empty images")
	ErrFailedScale      = errors.New("failed scale image")
	ErrUnsupportedType  = errors.New("unknown extension")
)
