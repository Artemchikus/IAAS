package postgres

import (
	"IAAS/internal/config"
	"IAAS/internal/models"
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
	config.Clusters = TestClusters(t)

	return config
}

func TestClusters(t *testing.T) []*models.Cluster {
	cluster := &models.Cluster{
		Location: "rus",
		URL:      "rus",
		Admin:    models.TestAdmin(t),
	}

	return []*models.Cluster{cluster}
}

func TestCluster(t *testing.T) *models.Cluster {
	return &models.Cluster{
		Location: "test",
		URL:      "test",
		Admin:    models.TestAdmin(t),
	}
}
