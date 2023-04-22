package models

import (
	"context"
	"testing"
)

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

func TestInitContext(t *testing.T) context.Context {
	return context.WithValue(context.Background(), CtxKeyRequestID, "test-initial-request")
}

func TestRequestContext(t *testing.T) context.Context {
	return context.WithValue(context.Background(), CtxKeyRequestID, "99999999-9999-9999-9999-999999999999")
}

func TestClusters(t *testing.T) []*Cluster {
	cluster := &Cluster{
		ID:       1,
		Location: "rus",
		URL:      "http://192.168.122.20",
		Admin:    TestClusterAdmin(t),
	}

	return []*Cluster{cluster}
}

func TestClusterAdmin(t *testing.T) *Account {
	return &Account{
		Email:     "adm@example.com",
		Name:      "admin",
		Password:  "openstack",
		ProjectID: "a0f8bb68921b4ddea63c52ef1b4e7e74",
	}
}
