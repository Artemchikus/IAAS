package api

import "errors"

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAutheticated          = errors.New("not authenticated")
	errEmailAlreadyExists       = errors.New("email is already registered")
	errWorngId                  = errors.New("invalid id given")
	errIncorrectRefreshToken    = errors.New("incorrect refresh token")
	errBadRequest               = errors.New("bad request")
	errInternalServerError      = errors.New("internal server error")
)
