package models

type Cluster struct {
	ID       int          `json:"id"`
	Location string       `json:"location"`
	Admin    *ClusterUser `json:"admin" toml:"cluster_admin"`
	URL      string       `json:"url"`
}

func NewCluster(adminName, adminEmail, adminPassword, adminDomainId, adminProjectId, url, location string) *Cluster {
	clusterAdmin := NewClusterUser(adminName, adminEmail, adminPassword, adminProjectId, adminDomainId, "")

	return &Cluster{
		Location: location,
		URL:      url,
		Admin:    clusterAdmin,
	}
}
