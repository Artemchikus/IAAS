package business

import (
	"IAAS/internal/models"
	"context"
)

type UserFetcher interface {
	FetchByID(context.Context, string) (*models.ClusterUser, error)
	Create(context.Context, *models.ClusterUser) error
	Delete(context.Context, string) error
	FetchAll(context.Context) ([]*models.ClusterUser, error)
}
type FlavorFetcher interface {
	FetchByID(context.Context, string) (*models.Flavor, error)
	Create(context.Context, *models.Flavor) error
	Delete(context.Context, string) error
	FetchAll(context.Context) ([]*models.Flavor, error)
}
type FloatingIPFetcher interface {
	FetchByID(context.Context, string) (*models.FloatingIp, error)
	Create(context.Context, *models.FloatingIp) error
	Delete(context.Context, string) error
	AddToPort(context.Context, string, string) error
	FetchAll(context.Context) ([]*models.FloatingIp, error)
}
type ImageFetcher interface {
	FetchByID(context.Context, string) (*models.Image, error)
	Create(context.Context, *models.Image) error
	Delete(context.Context, string) error
	FetchAll(context.Context) ([]*models.Image, error)
	Upload(context.Context, []byte, string) error
}
type NetworkFetcher interface {
	FetchByID(context.Context, string) (*models.Network, error)
	Create(context.Context, *models.Network) error
	Delete(context.Context, string) error
	FetchAll(context.Context) ([]*models.Network, error)
}
type PortFetcher interface {
	FetchByID(context.Context, string) (*models.Port, error)
	FetchByNetworkID(context.Context, string) ([]*models.Port, error)
	FetchByDeviceID(context.Context, string) ([]*models.Port, error)
	FetchAll(context.Context) ([]*models.Port, error)
}
type ProjectFetcher interface {
	FetchByID(context.Context, string) (*models.Project, error)
	Create(context.Context, *models.Project) error
	Delete(context.Context, string) error
	FetchAll(context.Context) ([]*models.Project, error)
}
type RoleFetcher interface {
	FetchByID(context.Context, string) (*models.Role, error)
	Create(context.Context, *models.Role) error
	Delete(context.Context, string) error
	FetchAll(context.Context) ([]*models.Role, error)
}
type RouterFetcher interface {
	FetchByID(context.Context, string) (*models.Router, error)
	Create(context.Context, *models.Router) error
	Delete(context.Context, string) error
	AddSubnet(context.Context, string, string) error
	RemoveSubnet(context.Context, string, string) error
	RemoveExternalGateway(context.Context, string) error
	FetchAll(context.Context) ([]*models.Router, error)
}
type SecurityGroupFetcher interface {
	FetchByID(context.Context, string) (*models.SecurityGroup, error)
	Create(context.Context, *models.SecurityGroup) error
	Delete(context.Context, string) error
	FetchAll(context.Context) ([]*models.SecurityGroup, error)
}
type SecurityRuleFetcher interface {
	FetchByID(context.Context, string) (*models.SecurityRule, error)
	Create(context.Context, *models.SecurityRule) error
	Delete(context.Context, string) error
	FetchAll(context.Context) ([]*models.SecurityRule, error)
}
type ServerFetcher interface {
	FetchByID(context.Context, string) (*models.Server, error)
	Create(context.Context, *models.Server) error
	Delete(context.Context, string) error
	Start(context.Context, string) error
	Stop(context.Context, string) error
	AttachVolume(context.Context, string, string) error
	FetchAll(context.Context) ([]*models.Server, error)
}
type KeyPairFetcher interface {
	FetchByID(context.Context, string) (*models.KeyPair, error)
	Create(context.Context, *models.KeyPair) error
	Delete(context.Context, string) error
	FetchAll(context.Context) ([]*models.KeyPair, error)
}
type SubnetFetcher interface {
	FetchByID(context.Context, string) (*models.Subnet, error)
	Create(context.Context, *models.Subnet) error
	Delete(context.Context, string) error
	FetchAll(context.Context) ([]*models.Subnet, error)
}
type TokenFetcher interface {
	Get(context.Context, *models.ClusterUser) (*models.Token, error)
	GetAdmin(context.Context) (*models.Token, error)
}
type VolumeFetcher interface {
	FetchByID(context.Context, string) (*models.Volume, error)
	Create(context.Context, *models.Volume) error
	Delete(context.Context, string) error
	FetchAll(context.Context) ([]*models.Volume, error)
}
