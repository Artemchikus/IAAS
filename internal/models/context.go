package models

type ctxKey int8

const (
	CtxKeyRequestID ctxKey = iota
	CtxKeyAccount
)
