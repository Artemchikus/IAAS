package business

import (
	"IAAS/internal/models"
	"context"
)

type UserFetcher interface {
	FetchByID(context.Context, string) (*models.ClusterUser, error)
	Create(context.Context, *models.ClusterUser) error
	Delete(context.Context, string) error
}
type FlavorFetcher interface {
	FetchByID(context.Context, string) (*models.Flavor, error)
	Create(context.Context, *models.Flavor) error
	Delete(context.Context, string) error
}
type FloatingIPFetcher interface {
	FetchByID(context.Context, string) (*models.FloatingIp, error)
	Create(context.Context, *models.FloatingIp) error
	Delete(context.Context, string) error
}
type ImageFetcher interface {
	FetchByID(context.Context, string) (*models.Image, error)
	Create(context.Context, *models.Image) error
	Delete(context.Context, string) error
}
type NetworkFetcher interface {
	FetchByID(context.Context, string) (*models.Network, error)
	Create(context.Context, *models.Network) error
	Delete(context.Context, string) error
}
type PortFetcher interface {
	FetchByID(context.Context, string) (*models.Port, error)
	FetchByNetworkID(context.Context, string) ([]*models.Port, error)
	FetchByRouterID(context.Context, string) ([]*models.Port, error)
}
type ProjectFetcher interface {
	FetchByID(context.Context, string) (*models.Project, error)
	Create(context.Context, *models.Project) error
	Delete(context.Context, string) error
}
type RoleFetcher interface {
	FetchByID(context.Context, string) (*models.Role, error)
	Create(context.Context, *models.Role) error
	Delete(context.Context, string) error
}
type RouterFetcher interface {
	FetchByID(context.Context, string) (*models.Router, error)
	Create(context.Context, *models.Router) error
	Delete(context.Context, string) error
}
type SecurityGroupFetcher interface {
	FetchByID(context.Context, string) (*models.SecurityGroup, error)
	Create(context.Context, *models.SecurityGroup) error
	Delete(context.Context, string) error
}
type SecurityRuleFetcher interface {
	FetchByID(context.Context, string) (*models.SecurityRule, error)
	Create(context.Context, *models.SecurityRule) error
	Delete(context.Context, string) error
}
type ServerFetcher interface {
	FetchByID(context.Context, string) (*models.Server, error)
	Create(context.Context, *models.Server) error
	Delete(context.Context, string) error
}
type KeyPairFetcher interface {
	FetchByID(context.Context, string) (*models.KeyPair, error)
	Create(context.Context, *models.KeyPair) error
	Delete(context.Context, string) error
}
type SubnetFetcher interface {
	FetchByID(context.Context, string) (*models.Subnet, error)
	Create(context.Context, *models.Subnet) error
	Delete(context.Context, string) error
}
type TokenFetcher interface {
	Get(context.Context, *models.ClusterUser) (*models.Token, error)
	GetAdmin(context.Context) (*models.Token, error)
}
type VolumeFetcher interface {
	FetchByID(context.Context, string) (*models.Volume, error)
	Create(context.Context, *models.Volume) error
	Delete(context.Context, string) error
}
type VolumeAttachmentFetcher interface {
	FetchByID(context.Context, string) (*models.VolumeAttachment, error)
	Create(context.Context, *models.VolumeAttachment) error
	Delete(context.Context, string) error
}
