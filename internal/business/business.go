package business

type Fetcher interface {
	Server() ServerFetcher
	Flavor() FlavorFetcher
	// Volume() VolumeFetcher
	// VolumeAttachment() VolumeAttachmentFetcher
	User() UserFetcher
	Token() TokenFetcher
	Image() ImageFetcher
	Subnet() SubnetFetcher
	Network() NetworkFetcher
	// SSHKey() SSHKeyFetcher
	SecurityGroup() SecurityGroupFetcher // *
	Router() RouterFetcher               // *
	Role() RoleFetcher                   // *
	FloatingIp() FloatingIPFetcher
	// Port() PortFetcher
	Project() ProjectFetcher
	SecurityRule() SecurityRuleFetcher
}
