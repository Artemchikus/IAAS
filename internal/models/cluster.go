package models

type Cluster struct {
	ID       int      `json:"id"`
	Location string   `json:"location"`
	Admin    *Account `json:"admin" toml:"cluster_admin"`
	URL      string   `json:"url"`
}

func NewCluster(adminName, adminEmail, adminPassword, url, location string) (*Cluster, error) {
	clusterAdmin, err := NewAccount(adminName, adminEmail, adminPassword)
	if err != nil {
		return nil, err
	}

	return &Cluster{
		Location: location,
		URL:      url,
		Admin:    clusterAdmin,
	}, nil
}
