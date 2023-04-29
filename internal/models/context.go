package models

type ctxKey int8

const (
	CtxKeyRequestID ctxKey = iota
	CtxKeyAccount
	CtxClusterUser
	CtxKeyClusterID
	CtxKeyToken
	CtxKeyProjectID
)
