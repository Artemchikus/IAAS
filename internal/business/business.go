package business

type Fetcher interface {
	Server() ServerFetcher
	Flavor() FlavorFetcher
	Volume() VolumeFetcher
	// VolumeAttachment() VolumeAttachmentFetcher -H "OpenStack-API-Version: volume 3.27"
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
}
