package business

import (
	"IAAS/internal/store"
	"context"
)

type Fetcher interface {
	Server() ServerFetcher
	Flavor() FlavorFetcher
	Volume() VolumeFetcher
	User() UserFetcher
	Token() TokenFetcher
	Image() ImageFetcher
	Subnet() SubnetFetcher
	Network() NetworkFetcher
	KeyPair() KeyPairFetcher
	SecurityGroup() SecurityGroupFetcher
	Router() RouterFetcher
	Role() RoleFetcher
	FloatingIp() FloatingIPFetcher
	// Port() PortFetcher
	Project() ProjectFetcher
	SecurityRule() SecurityRuleFetcher
	UpdateClusterMap(context.Context, store.Storage) error
}
