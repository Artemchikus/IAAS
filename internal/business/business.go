package business

type Fetcher interface {
	Server() ServerFetcher
	// Flavor() FlavorFetcher
	// Volume() VolumeFetcher
	// VolumeAttachment() VolumeAttachmentFetcher
	User() UserFetcher
	Token() TokenFetcher
	// Cluster() ClusterFetcher
	// Subnet() SubnetFetcher
	// Network() NetworkFetcher
	// SSHKey() SSHKeyFetcher
	// SecurityGroup() SecurityGroupFetcher
	// Router() RouterFetcher
	// Role() RoleFetcher
	// FloatingIP() FloatingIPFetcher
	// Port() PortFetcher
	// Project() ProjectFetcher
	// SecurityRuler() SecurityRulerFetcher
}
