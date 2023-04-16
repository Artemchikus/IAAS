package business

import (
	"IAAS/internal/models"
	"context"
)

type ServerFetcher interface {
	FetchByID(context.Context, int)
	Create(context.Context, int)
	Delete(context.Context, int)
	Update(context.Context, int)
}

type TokenFetcher interface {
	Get(context.Context, int, *models.Account) (*models.Token, error)
}

type FlavorFetcher interface {
	FetchByID(context.Context, int)
	Create(context.Context, int)
	Delete(context.Context, int)
	Update(context.Context, int)
}

type VolumeFetcher interface {
	FetchByID(context.Context, int)
	Create(context.Context, int)
	Delete(context.Context, int)
	Update(context.Context, int)
}

type VolumeAttachmentFetcher interface {
	FetchByID(context.Context, int)
	Create(context.Context, int)
	Delete(context.Context, int)
	Update(context.Context, int)
}

type UserFetcher interface {
	FetchByID(context.Context, int, string) (*models.Account, error)
	Create(context.Context, int, *models.Account) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}

type ClusterFetcher interface {
	FetchByID(context.Context, int)
	Create(context.Context, int)
	Delete(context.Context, int)
	Update(context.Context, int)
}

type SubnetFetcher interface {
	FetchByID(context.Context, int)
	Create(context.Context, int)
	Delete(context.Context, int)
	Update(context.Context, int)
}

type NetworkFetcher interface {
	FetchByID(context.Context, int)
	Create(context.Context, int)
	Delete(context.Context, int)
	Update(context.Context, int)
}

type SSHKeyFetcher interface {
	FetchByID(context.Context, int)
	Create(context.Context, int)
	Delete(context.Context, int)
	Update(context.Context, int)
}

type SecurityGroupFetcher interface {
	FetchByID(context.Context, int)
	Create(context.Context, int)
	Delete(context.Context, int)
	Update(context.Context, int)
}

type RouterFetcher interface {
	FetchByID(context.Context, int)
	Create(context.Context, int)
	Delete(context.Context, int)
	Update(context.Context, int)
}

type RoleFetcher interface {
	FetchByID(context.Context, int)
	Create(context.Context, int)
	Delete(context.Context, int)
	Update(context.Context, int)
}

type FloatingIPFetcher interface {
	FetchByID(context.Context, int)
	Create(context.Context, int)
	Delete(context.Context, int)
	Update(context.Context, int)
}

type PortFetcher interface {
	FetchByID(context.Context, int)
	Create(context.Context, int)
	Delete(context.Context, int)
	Update(context.Context, int)
}

type ProjectFetcher interface {
	FetchByID(context.Context, int, string) (*models.Project, error)
	Create(context.Context, int, *models.Project) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}

type SecurityRulerFetcher interface {
	FetchByID(context.Context, int)
	Create(context.Context, int)
	Delete(context.Context, int)
	Update(context.Context, int)
}
