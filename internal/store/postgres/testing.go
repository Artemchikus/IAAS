package postgres

import (
	"IAAS/internal/config"
	"IAAS/internal/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

func TestDB(t *testing.T, databaseURL string) (*sql.DB, func(...string)) {
	t.Helper()

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			db.Exec(fmt.Sprintf("TRUNCATE %s RESTART IDENTITY CASCADE", strings.Join(tables, ", ")))
		}

		db.Close()
	}
}

func TestConfig(t *testing.T) *config.ApiConfig {
	config := config.NewConfig()
	config.JwtKey = "secretkey"
	config.Admin = models.TestAdmin(t)
	config.Clusters = models.TestClusters(t)

	return config
}

func TestInitContext(t *testing.T) context.Context {
	return context.WithValue(context.Background(), models.CtxKeyRequestID, "test-initial-request")
}

func TestRequestContext(t *testing.T) context.Context {
	return context.WithValue(context.Background(), models.CtxKeyRequestID, "99999999-9999-9999-9999-999999999999")

}
