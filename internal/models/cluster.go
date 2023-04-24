package models

type Cluster struct {
	ID       int      `json:"id"`
	Location string   `json:"location"`
	Admin    *Account `json:"admin" toml:"cluster_admin"`
	URL      string   `json:"url"`
}

func NewCluster(adminName, adminEmail, adminPassword, url, location string) *Cluster {
	clusterAdmin := NewAccount(adminName, adminEmail, adminPassword, "admin")

	return &Cluster{
		Location: location,
		URL:      url,
		Admin:    clusterAdmin,
	}
}
