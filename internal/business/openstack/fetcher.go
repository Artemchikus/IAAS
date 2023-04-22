package openstack

import (
	"IAAS/internal/business"
	"IAAS/internal/config"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"context"
	"net/http"

	"go.uber.org/zap"
)

type Fetcher struct {
	logger               *zap.SugaredLogger
	client               *http.Client
	serverFetcher        *ServerFetcher
	userFetcher          *UserFetcher
	tokenFetcher         *TokenFetcher
	projectFetcher       *ProjectFetcher
	imageFetcher         *ImageFetcher
	flavorFetcher        *FlavorFetcher
	floatingIpFetcher    *FloatingIpFetcher
	networkFetcher       *NetworkFetcher
	subnetFetcher        *SubnetFetcher
	roleFetcher          *RoleFetcher
	routerFecther        *RouterFetcher
	securityGroupFetcher *SecurityGroupFetcher
	securityRuleFetcher  *SecurityRuleFetcher
	keyPairFetcher       *KeyPairFetcher
	volumeFetcher        *VolumeFetcher
	clusters             map[int]*models.Cluster
}

func NewFetcher(ctx context.Context, config *config.ApiConfig, store *postgres.Store) *Fetcher {
	log, _ := zap.NewProduction()
	defer log.Sync()
	sugar := log.Sugar()

	client := &http.Client{}

	clusters, err := store.Cluster().GetAll(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	clusterMap := make(map[int]*models.Cluster)

	for _, cluster := range clusters {
		clusterMap[cluster.ID] = cluster
	}

	return &Fetcher{
		logger:   sugar,
		client:   client,
		clusters: clusterMap,
	}
}

func (f *Fetcher) Server() business.ServerFetcher {
	if f.serverFetcher != nil {
		return f.serverFetcher
	}

	f.serverFetcher = &ServerFetcher{
		fetcher: f,
	}

	return f.serverFetcher
}

func (f *Fetcher) User() business.UserFetcher {
	if f.userFetcher != nil {
		return f.userFetcher
	}

	f.userFetcher = &UserFetcher{
		fetcher: f,
	}

	return f.userFetcher
}

func (f *Fetcher) Token() business.TokenFetcher {
	if f.tokenFetcher != nil {
		return f.tokenFetcher
	}

	f.tokenFetcher = &TokenFetcher{
		fetcher: f,
	}

	return f.tokenFetcher
}

func (f *Fetcher) Project() business.ProjectFetcher {
	if f.projectFetcher != nil {
		return f.projectFetcher
	}

	f.projectFetcher = &ProjectFetcher{
		fetcher: f,
	}

	return f.projectFetcher
}

func (f *Fetcher) Image() business.ImageFetcher {
	if f.imageFetcher != nil {
		return f.imageFetcher
	}

	f.imageFetcher = &ImageFetcher{
		fetcher: f,
	}

	return f.imageFetcher
}

func (f *Fetcher) Flavor() business.FlavorFetcher {
	if f.flavorFetcher != nil {
		return f.flavorFetcher
	}

	f.flavorFetcher = &FlavorFetcher{
		fetcher: f,
	}

	return f.flavorFetcher
}

func (f *Fetcher) FloatingIp() business.FloatingIPFetcher {
	if f.floatingIpFetcher != nil {
		return f.floatingIpFetcher
	}

	f.floatingIpFetcher = &FloatingIpFetcher{
		fetcher: f,
	}

	return f.floatingIpFetcher
}

func (f *Fetcher) Network() business.NetworkFetcher {
	if f.networkFetcher != nil {
		return f.networkFetcher
	}

	f.networkFetcher = &NetworkFetcher{
		fetcher: f,
	}

	return f.networkFetcher
}

func (f *Fetcher) Subnet() business.SubnetFetcher {
	if f.subnetFetcher != nil {
		return f.subnetFetcher
	}

	f.subnetFetcher = &SubnetFetcher{
		fetcher: f,
	}

	return f.subnetFetcher
}

func (f *Fetcher) Role() business.RoleFetcher {
	if f.roleFetcher != nil {
		return f.roleFetcher
	}

	f.roleFetcher = &RoleFetcher{
		fetcher: f,
	}

	return f.roleFetcher
}

func (f *Fetcher) Router() business.RouterFetcher {
	if f.routerFecther != nil {
		return f.routerFecther
	}

	f.routerFecther = &RouterFetcher{
		fetcher: f,
	}

	return f.routerFecther
}

func (f *Fetcher) SecurityGroup() business.SecurityGroupFetcher {
	if f.securityGroupFetcher != nil {
		return f.securityGroupFetcher
	}

	f.securityGroupFetcher = &SecurityGroupFetcher{
		fetcher: f,
	}

	return f.securityGroupFetcher
}

func (f *Fetcher) SecurityRule() business.SecurityRuleFetcher {
	if f.securityRuleFetcher != nil {
		return f.securityRuleFetcher
	}

	f.securityRuleFetcher = &SecurityRuleFetcher{
		fetcher: f,
	}

	return f.securityRuleFetcher
}

func (f *Fetcher) KeyPair() business.KeyPairFetcher {
	if f.keyPairFetcher != nil {
		return f.keyPairFetcher
	}

	f.keyPairFetcher = &KeyPairFetcher{
		fetcher: f,
	}

	return f.keyPairFetcher
}

func (f *Fetcher) Volume() business.VolumeFetcher {
	if f.volumeFetcher != nil {
		return f.volumeFetcher
	}

	f.volumeFetcher = &VolumeFetcher{
		fetcher: f,
	}

	return f.volumeFetcher
}

func getClusterIDFromContext(ctx context.Context) int {
	return ctx.Value(models.CtxKeyClusterID).(int)
}

func getTokenFromContext(ctx context.Context) *models.Token {
	return ctx.Value(models.CtxKeyClusterID).(*models.Token)
}
