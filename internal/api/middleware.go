package api

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sugar := s.logger.With(
			"remoute_addr", r.RemoteAddr,
			"request_id", r.Context().Value(ctxKeyRequestID),
		)
		sugar.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		sugar.Infof("complited with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Since(start))
	})
}

func (s *server) authenticateAccount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("x-jwt-token")

		secret, err := s.store.Secret().GetByType("jwt")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, errInternalServerError)
			return
		}

		userID, err := getId(r)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errBadRequest)
			return
		}
		account, err := s.store.Account().FindByID(userID)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAutheticated)
			return
		}

		if err := account.ValidateJWTToken(tokenStr, secret.Value); err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAutheticated)
			return
		}

		next.ServeHTTP(w, r)
	})
}
