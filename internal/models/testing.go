package models

import "testing"

func TestAccount(t *testing.T) *Account {
	return &Account{
		Email:    "test@example.com",
		Name:     "test",
		Password: "password",
	}
}

func TestAdmin(t *testing.T) *Account {
	return &Account{
		Email:    "test_adm@example.com",
		Name:     "test_adm",
		Password: "password",
	}
}

func TestClusters(t *testing.T) []*Cluster {
	cluster := &Cluster{
		Location: "rus",
		URL:      "rus",
		Admin:    TestAdmin(t),
	}

	return []*Cluster{cluster}
}

func TestCluster(t *testing.T) *Cluster {
	return &Cluster{
		Location: "test",
		URL:      "test",
		Admin:    TestAdmin(t),
	}
}
