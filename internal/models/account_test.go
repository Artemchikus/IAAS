package models_test

import (
	"IAAS/internal/models"
	"fmt"
	"testing"
)

func TestNewAccount(t *testing.T) {
	acc := models.NewAccount("a", "b", "c", "user")
	fmt.Printf("%+v\n", acc)
}
