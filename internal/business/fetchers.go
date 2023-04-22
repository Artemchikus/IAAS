package business

import (
	"IAAS/internal/models"
	"context"
)

type UserFetcher interface {
	FetchByID(context.Context, string) (*models.Account, error)
	Create(context.Context, *models.Account) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
type FlavorFetcher interface {
	FetchByID(context.Context, string) (*models.Flavor, error)
	Create(context.Context, *models.Flavor) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
type FloatingIPFetcher interface {
	FetchByID(context.Context, string) (*models.FloatingIp, error)
	Create(context.Context, *models.FloatingIp) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
type ImageFetcher interface {
	FetchByID(context.Context, string) (*models.Image, error)
	Create(context.Context, *models.Image) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
type NetworkFetcher interface {
	FetchByID(context.Context, string) (*models.Network, error)
	Create(context.Context, *models.Network) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
type PortFetcher interface {
	FetchByID(context.Context, string) (*models.Port, error)
	Create(context.Context, *models.Port) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
type ProjectFetcher interface {
	FetchByID(context.Context, string) (*models.Project, error)
	Create(context.Context, *models.Project) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
type RoleFetcher interface {
	FetchByID(context.Context, string) (*models.Role, error)
	Create(context.Context, *models.Role) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
type RouterFetcher interface {
	FetchByID(context.Context, string) (*models.Router, error)
	Create(context.Context, *models.Router) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
type SecurityGroupFetcher interface {
	FetchByID(context.Context, string) (*models.SecurityGroup, error)
	Create(context.Context, *models.SecurityGroup) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
type SecurityRuleFetcher interface {
	FetchByID(context.Context, string) (*models.SecurityRule, error)
	Create(context.Context, *models.SecurityRule) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
type ServerFetcher interface {
	FetchByID(context.Context, string) (*models.Server, error)
	Create(context.Context, *models.Server) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
type KeyPairFetcher interface {
	FetchByID(context.Context, string) (*models.KeyPair, error)
	Create(context.Context, *models.KeyPair) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
type SubnetFetcher interface {
	FetchByID(context.Context, string) (*models.Subnet, error)
	Create(context.Context, *models.Subnet) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
type TokenFetcher interface {
	Get(context.Context, *models.Account) (*models.Token, error)
}
type VolumeFetcher interface {
	FetchByID(context.Context, string) (*models.Volume, error)
	Create(context.Context, *models.Volume) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
type VolumeAttachmentFetcher interface {
	FetchByID(context.Context, string) (*models.VolumeAttachment, error)
	Create(context.Context, *models.VolumeAttachment) error
	Delete(context.Context, string) error
	Update(context.Context, string)
}
