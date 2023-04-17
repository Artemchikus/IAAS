package business

import (
	"IAAS/internal/models"
	"context"
)

type UserFetcher interface {
	FetchByID(context.Context, int, string) (*models.Account, error)
	Create(context.Context, int, *models.Account) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
type FlavorFetcher interface {
	FetchByID(context.Context, int, string) (*models.Flavor, error)
	Create(context.Context, int, *models.Flavor) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
type FloatingIPFetcher interface {
	FetchByID(context.Context, int, string) (*models.FloatingIp, error)
	Create(context.Context, int, *models.FloatingIp) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
type ImageFetcher interface {
	FetchByID(context.Context, int, string) (*models.Image, error)
	Create(context.Context, int, *models.Image) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
type NetworkFetcher interface {
	FetchByID(context.Context, int, string) (*models.Network, error)
	Create(context.Context, int, *models.Network) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
type PortFetcher interface {
	FetchByID(context.Context, int, string) (*models.Port, error)
	Create(context.Context, int, *models.Port) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
type ProjectFetcher interface {
	FetchByID(context.Context, int, string) (*models.Project, error)
	Create(context.Context, int, *models.Project) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
type RoleFetcher interface {
	FetchByID(context.Context, int, string) (*models.Role, error)
	Create(context.Context, int, *models.Role) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
type RouterFetcher interface {
	FetchByID(context.Context, int, string) (*models.Router, error)
	Create(context.Context, int, *models.Router) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
type SecurityGroupFetcher interface {
	FetchByID(context.Context, int, string) (*models.SecurityGroup, error)
	Create(context.Context, int, *models.SecurityGroup) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
type SecurityRuleFetcher interface {
	FetchByID(context.Context, int, string) (*models.SecurityRule, error)
	Create(context.Context, int, *models.SecurityRule) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
type ServerFetcher interface {
	FetchByID(context.Context, int, string) (*models.Server, error)
	Create(context.Context, int, *models.Server) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
type SSHKeyFetcher interface {
	FetchByID(context.Context, int, string) (*models.SSHKey, error)
	Create(context.Context, int, *models.SSHKey) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
type SubnetFetcher interface {
	FetchByID(context.Context, int, string) (*models.Subnet, error)
	Create(context.Context, int, *models.Subnet) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
type TokenFetcher interface {
	Get(context.Context, int, *models.Account) (*models.Token, error)
}
type VolumeFetcher interface {
	FetchByID(context.Context, int, string) (*models.Volume, error)
	Create(context.Context, int, *models.Volume) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
type VolumeAttachmentFetcher interface {
	FetchByID(context.Context, int, string) (*models.VolumeAttachment, error)
	Create(context.Context, int, *models.VolumeAttachment) error
	Delete(context.Context, int, string) error
	Update(context.Context, int)
}
