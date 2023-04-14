package models

type Cluster struct {
	ID       int      `json:"id"`
	Location string   `json:"location"`
	Admin    *Account `json:"admin" toml:"cluster_admin"`
	URL      string   `json:"url"`
}
