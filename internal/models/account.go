package models

import (
	"errors"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID                int       `json:"id"`
	ProjectID         string    `json:"project_id" toml:"project_id"`
	Name              string    `json:"name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Email             string    `json:"email"`
	EncryptedPassword string    `json:"encrypted_password"`
	Password          string    `json:"password,omitempty"`
	RefreshToken      string    `json:"refresh_token"`
	Role              string    `json:"role"`
}

func NewAccount(name, email, password, role string) *Account {
	return &Account{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Email:     email,
		UpdatedAt: time.Now().UTC(),
		Password:  password,
		Role:      role,
	}
}

func (a *Account) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, validation.By(requiredIf(a.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

func (a *Account) BeforeCreate() error {
	enc, err := encryptSrting(a.Password)
	if err != nil {
		return err
	}

	a.EncryptedPassword = enc
	return nil
}

func (a *Account) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(password)) == nil
}

func (a *Account) ValidateJWTToken(tokenStr, secret string) error {
	token, err := parseToken(tokenStr, secret)
	if err != nil {
		return errors.New("incorrect token")
	}

	if !token.Valid {
		return errors.New("not authenticated")
	}

	claims := token.Claims.(jwt.MapClaims)
	if float64(time.Now().UTC().Unix()) >= claims["ExpiresAt"].(float64) {
		return errors.New("token expired")
	}

	if a.Email != claims["AccountEmail"] {
		return errors.New("not authenticated")
	}

	return nil
}

func parseToken(tokenStr, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func encryptSrting(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// {
// 	"default_project_id": "2e83a8e5362247f79c8d86980ab2a216",
// 	"domain_id": "default",
// 	"enabled": true,
// 	"id": "d8b5c61b6e594b218bca0f9e530a9f95",
// 	"name": "octavia",
// 	"options": {},
// 	"password_expires_at": null
//   }
