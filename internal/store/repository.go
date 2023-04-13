package store

import (
	"IAAS/internal/models"
	"context"
	"time"
)

type AccountRepository interface {
	Create(context.Context, *models.Account) error
	Delete(context.Context, int) error
	Update(context.Context, *models.Account) error
	FindByID(context.Context, int) (*models.Account, error)
	FindByEmail(context.Context, string) (*models.Account, error)
	Init(context.Context) error
	GetAll(context.Context) ([]*models.Account, error)
	UpdateRefreshToken(context.Context, string, string, time.Time) error
}

type SecretRepository interface {
	GetByType(context.Context, string) (*models.Secret, error)
	Init(context.Context, *models.Secret) error
}
