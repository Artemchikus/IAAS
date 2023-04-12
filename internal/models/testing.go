package models

import "testing"

func TestAccount(t *testing.T) *Account {
	return &Account{
		Email:    "test@example.com",
		Name:     "test",
		Password: "password",
	}
}
