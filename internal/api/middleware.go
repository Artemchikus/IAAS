package api

import (
	"IAAS/internal/models"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), models.CtxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sugar := s.logger.With(
			"remoute_addr", r.RemoteAddr,
			"request_id", r.Context().Value(models.CtxKeyRequestID),
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

		secret, err := s.store.Secret().GetByType(r.Context(), "jwt")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		userID, err := getId(r)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		account, err := s.store.Account().FindByID(r.Context(), userID)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAutheticated)
			return
		}

		if err := account.ValidateJWTToken(tokenStr, secret.Value); err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
