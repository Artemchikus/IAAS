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

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) error {
	req := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	acc, err := s.store.Account().FindByEmail(req.Email)
	if err != nil {
		return err
	}

	if !acc.ValidatePassword(req.Password) {
		return errIncorrectEmailOrPassword
	}

	secret, err := s.store.Secret().GetByType("jwt")
	if err != nil {
		return err
	}

	jwtToken, err := createJWT(acc, secret.Value)
	if err != nil {
		return err
	}

	refreshToken, err := createRefreshToken(acc, secret.Value)
	if err != nil {
		return err
	}

	if err := s.store.Account().UpdateRefreshToken(acc.RefreshToken, refreshToken, time.Now()); err != nil {
		return err
	}

	resp := &LoginResponse{
		JWT:     jwtToken,
		Refresh: refreshToken,
	}

	return wirteJSON(w, http.StatusOK, resp)
}

func (s *server) handleGetAllAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.Account().GetAll()
	if err != nil {
		return err
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

	return wirteJSON(w, http.StatusOK, resp)
}

func (s *server) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}
	account, err := s.store.Account().FindByID(id)
	if err != nil {
		return err
	}

	resp := GetAccountResponse{
		ID:    account.ID,
		Name:  account.Name,
		Email: account.Email,
	}

	return wirteJSON(w, http.StatusOK, resp)
}

func (s *server) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account, err := models.NewAccount(req.Name, req.Email, req.Password)
	if err != nil {
		return err
	}

	if _, err := s.store.Account().FindByEmail(account.Email); err == nil {
		return errEmailAlreadyExists
	}

	if err := s.store.Account().Create(account); err != nil {
		return err
	}

	resp := CreateAccountResponse{
		ID:    account.ID,
		Email: account.Email,
		Name:  account.Name,
	}

	return wirteJSON(w, http.StatusOK, resp)
}

func (s *server) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}
	if err := s.store.Account().Delete(id); err != nil {
		return err
	}

	resp := DeleteAccountResponse{
		Deleted: id,
	}

	return wirteJSON(w, http.StatusOK, resp)
}

func (s *server) handleRefreshToken(w http.ResponseWriter, r *http.Request) error {
	refreshTokenStr := r.Header.Get("x-refresh-token")

	secret, err := s.store.Secret().GetByType("jwt")
	if err != nil {
		return err
	}

	acc, err := getAccFromRefreshToken(s.store, refreshTokenStr, secret.Value)
	if err != nil {
		return err
	}

	if acc.RefreshToken != refreshTokenStr {
		return errNotAutheticated
	}

	jwtToken, err := createJWT(acc, secret.Value)
	if err != nil {
		return err
	}

	refreshToken, err := createRefreshToken(acc, secret.Value)
	if err != nil {
		return err
	}

	if err := s.store.Account().UpdateRefreshToken(acc.RefreshToken, refreshToken, time.Now()); err != nil {
		return err
	}

	resp := &RefreshTokenResponse{
		Refresh: refreshToken,
		JWT:     jwtToken,
	}

	return wirteJSON(w, http.StatusOK, resp)
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
		return nil, errors.New("not authenticated")
	}

	return acc, nil
}
