package transport

import (
    "aws_rekognition_demo/internal/domain/auth"
    "errors"
    "net/http"
)

func TranslateError(err error) Error {
    if errors.Is(err, ErrBadRouting) {
        return Error{
            Status: http.StatusNotFound,
            Response: responseWithoutData{
                Code:    100,
                Message: err.Error(),
            },
        }
    } else if errors.Is(err, ErrBadRequest) {
        return Error{
            Status: http.StatusBadRequest,
            Response: responseWithoutData{
                Code:    101,
                Message: err.Error(),
            },
        }
    } else if errors.Is(err, auth.ErrUserNotExist) {
        return Error{
            Status: http.StatusNotFound,
            Response: responseWithoutData{
                Code:    102,
                Message: err.Error(),
            },
        }
    } else if errors.Is(err, auth.ErrInvalidAuthorizationHeader) {
        return Error{
            Status: http.StatusBadRequest,
            Response: responseWithoutData{
                Code:    103,
                Message: err.Error(),
            },
        }
    } else if errors.Is(err, auth.ErrInvalidAccessToken) {
        return Error{
            Status: http.StatusUnauthorized,
            Response: responseWithoutData{
                Code:    104,
                Message: err.Error(),
            },
        }
    } else if errors.Is(err, auth.ErrInvalidRefreshToken) {
        return Error{
            Status: http.StatusUnauthorized,
            Response: responseWithoutData{
                Code:    105,
                Message: err.Error(),
            },
        }
    } else if errors.Is(err, auth.ErrFailedCreateAccessToken) {
        return Error{
            Status: http.StatusInternalServerError,
            Response: responseWithoutData{
                Code:    106,
                Message: err.Error(),
            },
        }
    } else if errors.Is(err, auth.ErrFailedUpdateAccessToken) {
        return Error{
            Status: http.StatusInternalServerError,
            Response: responseWithoutData{
                Code:    107,
                Message: err.Error(),
            },
        }
    } else if errors.Is(err, auth.ErrFailedSaveUser) {
        return Error{
            Status: http.StatusInternalServerError,
            Response: responseWithoutData{
                Code:    108,
                Message: err.Error(),
            },
        }
    } else if errors.Is(err, auth.ErrFailedCreateUserID) {
        return Error{
            Status: http.StatusInternalServerError,
            Response: responseWithoutData{
                Code:    109,
                Message: err.Error(),
            },
        }
    } else if errors.Is(err, auth.ErrDuplicateUser) {
        return Error{
            Status: http.StatusConflict,
            Response: responseWithoutData{
                Code:    110,
                Message: err.Error(),
            },
        }
    } else if errors.Is(err, auth.ErrWrongPassword) {
        return Error{
            Status: http.StatusUnauthorized,
            Response: responseWithoutData{
                Code:    111,
                Message: err.Error(),
            },
        }
    }
    return Error{
        Status: http.StatusInternalServerError,
        Response: responseWithoutData{
            Code:    100,
            Message: "unexpected error",
        },
    }
}