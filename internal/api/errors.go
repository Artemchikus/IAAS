package api

import "errors"

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAutheticated          = errors.New("not authenticated")
	errEmailAlreadyExists       = errors.New("email is already registered")
	errIncorrectToken           = errors.New("incorrect token")
	errBadRequest               = errors.New("bad request")
	errInternalServerError      = errors.New("internal server error")
	errTokenExpired             = errors.New("token expired")
	errUnexpectedSigningMethod  = errors.New("unexpected signature method")
	errNoAdmin                  = errors.New("admin rights required")
	errLocationAlreadyExists    = errors.New("cluster location already exists")
	errAlreadyInCluster         = errors.New("already registered in cluster")
)
