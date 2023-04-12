package models_test

import (
	"IAAS/internal/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	acc, err := models.NewAccount("a", "b", "c")
	assert.Nil(t, err)
	fmt.Printf("%+v\n", acc)
}
