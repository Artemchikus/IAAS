package api

import (
	"IAAS/internal/models"
	"IAAS/internal/store"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {
	req := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	acc, err := s.store.Account().FindByEmail(req.Email)
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
		return
	}

	if !acc.ValidatePassword(req.Password) {
		s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
		return
	}

	secret, err := s.store.Secret().GetByType("jwt")
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, errInternalServerError)
		return
	}

	jwtToken, err := createJWT(acc, secret.Value)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, errInternalServerError)
		return
	}

	refreshToken, err := createRefreshToken(acc, secret.Value)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, errInternalServerError)
		return
	}

	if err := s.store.Account().UpdateRefreshToken(acc.RefreshToken, refreshToken, time.Now()); err != nil {
		s.error(w, r, http.StatusInternalServerError, errInternalServerError)
		return
	}

	resp := &LoginResponse{
		JWT:     jwtToken,
		Refresh: refreshToken,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.store.Account().GetAll()
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, errInternalServerError)
		return
	}

	resp := make(GetAllAccountsResponse, 0)

	for _, acc := range accounts {
		r := GetAccountResponse{
			ID:    acc.ID,
			Email: acc.Email,
			Name:  acc.Name,
		}
		resp = append(resp, &r)
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetAccountByID(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}
	account, err := s.store.Account().FindByID(id)
	if err != nil {
		s.error(w, r, http.StatusNotFound, errWorngId)
		return
	}

	resp := GetAccountResponse{
		ID:    account.ID,
		Name:  account.Name,
		Email: account.Email,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	account, err := models.NewAccount(req.Name, req.Email, req.Password)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, errInternalServerError)
		return
	}

	if _, err := s.store.Account().FindByEmail(account.Email); err == nil {
		s.error(w, r, http.StatusConflict, errEmailAlreadyExists)
		return
	}

	if err := account.Validate(); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := account.BeforeCreate(); err != nil {
		s.error(w, r, http.StatusInternalServerError, errInternalServerError)
		return
	}

	if err := s.store.Account().Create(account); err != nil {
		s.error(w, r, http.StatusInternalServerError, errInternalServerError)
		return
	}

	resp := CreateAccountResponse{
		ID:    account.ID,
		Email: account.Email,
		Name:  account.Name,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}
	if err := s.store.Account().Delete(id); err != nil {
		s.error(w, r, http.StatusInternalServerError, errInternalServerError)
		return
	}

	resp := DeleteAccountResponse{
		Deleted: id,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenStr := r.Header.Get("x-refresh-token")

	secret, err := s.store.Secret().GetByType("jwt")
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, errInternalServerError)
		return
	}

	acc, err := getAccFromRefreshToken(s.store, refreshTokenStr, secret.Value)
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, errIncorrectRefreshToken)
		return
	}

	if acc.RefreshToken != refreshTokenStr {
		s.error(w, r, http.StatusUnauthorized, errIncorrectRefreshToken)
		return
	}

	jwtToken, err := createJWT(acc, secret.Value)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, errInternalServerError)
		return
	}

	refreshToken, err := createRefreshToken(acc, secret.Value)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, errInternalServerError)
		return
	}

	if err := s.store.Account().UpdateRefreshToken(acc.RefreshToken, refreshToken, time.Now()); err != nil {
		s.error(w, r, http.StatusInternalServerError, errInternalServerError)
		return
	}

	resp := &RefreshTokenResponse{
		Refresh: refreshToken,
		JWT:     jwtToken,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func createJWT(account *models.Account, secret string) (string, error) {
	claims := &jwt.MapClaims{
		"ExpiresAt":    time.Now().Add(time.Minute * 30).Unix(),
		"AccountEmail": account.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func createRefreshToken(account *models.Account, secret string) (string, error) {
	claims := &jwt.MapClaims{
		"ExpiresAt":    time.Now().Add(time.Hour * 24).Unix(),
		"AccountEmail": account.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func getAccFromRefreshToken(s store.Storage, tokenStr, secret string) (*models.Account, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("not authenticated")
	}

	claims := token.Claims.(jwt.MapClaims)
	email := claims["AccountEmail"].(string)

	acc, err := s.Account().FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
