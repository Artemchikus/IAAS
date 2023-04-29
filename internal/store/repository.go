package store

import (
	"IAAS/internal/models"
	"context"
)

type AccountRepository interface {
	Create(context.Context, *models.Account) error
	Delete(context.Context, int) error
	Update(context.Context, *models.Account) error
	FindByID(context.Context, int) (*models.Account, error)
	FindByEmail(context.Context, string) (*models.Account, error)
	Init(context.Context, *models.Account) error // TODO add support for multiple accounts
	GetAll(context.Context) ([]*models.Account, error)
	UpdateRefreshToken(context.Context, string, string) error
}

type SecretRepository interface {
	FindByType(context.Context, string) (*models.Secret, error)
	Init(context.Context, *models.Secret) error // TODO add support for multiple secrets
}

type ClusterRepository interface {
	FindByID(context.Context, int) (*models.Cluster, error)
	FindByLocation(context.Context, string) (*models.Cluster, error)
	Create(context.Context, *models.Cluster) error
	Delete(context.Context, int) error
	Update(context.Context, *models.Cluster) error
	Init(context.Context, []*models.Cluster) error
	GetAll(context.Context) ([]*models.Cluster, error)
}

type ClusterUserRepository interface {
	Create(context.Context, *models.ClusterUser) error
	Delete(context.Context, string) error
	Update(context.Context, *models.ClusterUser) error
	FindByID(context.Context, string) (*models.ClusterUser, error)
	FindByEmailAndClusterID(context.Context, string, int) (*models.ClusterUser, error)
	FindByClusterID(context.Context, int) ([]*models.ClusterUser, error)
	FindByAccountID(context.Context, int) ([]*models.ClusterUser, error)
	Init(context.Context, []*models.Cluster) error
	GetAll(context.Context) ([]*models.ClusterUser, error)
}
