package transport

import (
	"errors"
	"github.com/col3name/images-search/pkg/app/dropbox"
	"github.com/col3name/images-search/pkg/domain/auth"
	"github.com/col3name/images-search/pkg/domain/pictures"
	"net/http"
)

var ErrUnexpected = errors.New("unexpected error")

func TranslateError(err error) *Error {
	if errorIs(err, ErrBadRouting) {
		return NewError(http.StatusNotFound, 100, err)
	} else if errorIs(err, ErrBadRequest) {
		return NewError(http.StatusBadRequest, 101, err)
	} else {
		e, isAuthError := tryTranslateAuthError(err)
		if isAuthError {
			return e
		} else if errorIs(err, pictures.ErrNotFound) {
			return NewError(http.StatusNotFound, 112, err)
		} else if errorIs(err, pictures.ErrEmptyImages) {
			return NewError(http.StatusNotFound, 113, err)
		} else if errorIs(err, pictures.ErrUnsupportedType) {
			return NewError(http.StatusBadRequest, 114, err)
		} else if errorIs(err, pictures.ErrFailedScale) {
			return NewError(http.StatusBadRequest, 114, err)
		} else {
			switch err.(type) {
			case *dropbox.ErrFailedGetListDropbox:
				return NewError(http.StatusBadRequest, 115, err)
			case *dropbox.ErrFailedDownloadDropbox:
				return NewError(http.StatusBadRequest, 116, err)
			}
		}
	}

	return NewError(http.StatusInternalServerError, 100, ErrUnexpected)
}

func errorIs(err, target error) bool {
	return errors.Is(err, target)
}

func tryTranslateAuthError(err error) (*Error, bool) {
	if errorIs(err, auth.ErrUserNotExist) {
		return NewError(http.StatusNotFound, 102, err), true
	} else if errorIs(err, auth.ErrInvalidAuthorizationHeader) {
		return NewError(http.StatusBadRequest, 103, err), true
	} else if errorIs(err, auth.ErrInvalidAccessToken) {
		return NewError(http.StatusUnauthorized, 104, err), true
	} else if errorIs(err, auth.ErrInvalidRefreshToken) {
		return NewError(http.StatusUnauthorized, 105, err), true
	} else if errorIs(err, auth.ErrFailedCreateAccessToken) {
		return NewError(http.StatusInternalServerError, 106, err), true
	} else if errorIs(err, auth.ErrFailedUpdateAccessToken) {
		return NewError(http.StatusInternalServerError, 107, err), true
	} else if errorIs(err, auth.ErrFailedSaveUser) {
		return NewError(http.StatusInternalServerError, 108, err), true
	} else if errorIs(err, auth.ErrFailedCreateUserID) {
		return NewError(http.StatusInternalServerError, 108, err), true
	} else if errorIs(err, auth.ErrDuplicateUser) {
		return NewError(http.StatusConflict, 110, err), true
	} else if errorIs(err, auth.ErrWrongPassword) {
		return NewError(http.StatusUnauthorized, 111, err), true
	}
	return nil, false
}
