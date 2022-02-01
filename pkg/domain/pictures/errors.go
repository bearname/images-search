package pictures

import "errors"

var ErrCannotAssignNull = errors.New("cannot assign NULL to *string")
var ErrEmptyImages = errors.New("empty images")
