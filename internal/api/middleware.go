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

		start := time.Now().UTC()
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

		secret, err := s.store.Secret().FindByType(r.Context(), "jwt")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		acc, err := getAccFromToken(r.Context(), s.store, tokenStr, secret.Value)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), models.CtxKeyAccount, acc)))
	})
}

func (s *server) isAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(models.CtxKeyAccount).(*models.Account).Role != "admin" {
			s.error(w, r, http.StatusUnauthorized, errNoAdmin)
			return
		}

		next.ServeHTTP(w, r)
	})
}
