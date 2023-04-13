package store

import (
	"IAAS/internal/models"
	"time"
)

type AccountRepository interface {
	Create( *models.Account) error
	Delete(int) error
	Update(*models.Account) error
	FindByID(int) (*models.Account, error)
	FindByEmail(string) (*models.Account, error)
	Init() error
	GetAll() ([]*models.Account, error)
	UpdateRefreshToken(string, string, time.Time) error
}

type SecretRepository interface {
	GetByType(string) (*models.Secret, error)
	Init(*models.Secret) error
}
