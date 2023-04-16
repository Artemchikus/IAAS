package postgres

import (
	"IAAS/internal/models"
	"IAAS/internal/store"
	"context"
	"database/sql"
	"errors"
	"time"
)

type ClusterRepository struct {
	store *Store
}

func (r *ClusterRepository) FindByID(ctx context.Context, id int) (*models.Cluster, error) {
	defer r.logging(ctx, "FIND BY id")()

	rows, err := r.store.db.Query("SELECT * FROM cluster WHERE id = $1", id)
	if err != nil {
		return nil, store.ErrRecordNotFound
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoCluster(rows)
	}

	return nil, store.ErrRecordNotFound
}

func (r *ClusterRepository) FindByLocation(ctx context.Context, location string) (*models.Cluster, error) {
	defer r.logging(ctx, "FIND BY location")()

	rows, err := r.store.db.Query("SELECT * FROM cluster WHERE location = $1", location)
	if err != nil {
		return nil, store.ErrRecordNotFound
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoCluster(rows)
	}

	return nil, store.ErrRecordNotFound
}

func (r *ClusterRepository) Create(ctx context.Context, cluster *models.Cluster) error {
	defer r.logging(ctx, "CREATE")()

	query := `
	INSERT INTO cluster 
	(location, url, admin_name, admin_email, admin_password)  
	values ($1, $2, $3, $4, $5) RETURNING *`

	row, err := r.store.db.Query(
		query,
		cluster.Location,
		cluster.URL,
		cluster.Admin.Name,
		cluster.Admin.Email,
		cluster.Admin.Password)
	if err != nil {
		return err
	}
	defer row.Close()

	for row.Next() {
		clas, err := scanIntoCluster(row)
		if err != nil {
			return err
		}
		cluster.ID = clas.ID
	}

	return nil
}

func (r *ClusterRepository) Delete(ctx context.Context, id int) error {
	defer r.logging(ctx, "DELETE")()

	row, err := r.store.db.Query("DELETE FROM cluster WHERE id = $1", id)
	if err != nil {
		return store.ErrRecordNotFound
	}
	defer row.Close()

	return nil
}

func (r *ClusterRepository) Update(ctx context.Context, cluster *models.Cluster) error {
	return nil
}

func (r *ClusterRepository) Init(ctx context.Context, clusters []*models.Cluster) error {
	defer r.logging(ctx, "INIT")()

	if err := r.createClusterTable(ctx); err != nil {
		return err
	}

	for _, cluster := range clusters {
		if err := r.addCluster(ctx, cluster); err != nil {
			return err
		}
	}

	return nil
}

func (r *ClusterRepository) GetAll(ctx context.Context) ([]*models.Cluster, error) {
	defer r.logging(ctx, "GET ALL")()

	rows, err := r.store.db.Query("SELECT * FROM cluster")
	if err != nil {
		return nil, store.ErrRecordNotFound
	}
	defer rows.Close()

	clusters := []*models.Cluster{}

	for rows.Next() {
		cluster, err := scanIntoCluster(rows)
		if err != nil {
			return nil, store.ErrRecordNotFound
		}
		clusters = append(clusters, cluster)
	}

	if len(clusters) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return clusters, nil
}

func scanIntoCluster(rows *sql.Rows) (*models.Cluster, error) {
	admin := new(models.Account)
	cluster := new(models.Cluster)
	cluster.Admin = admin

	if err := rows.Scan(
		&cluster.ID,
		&cluster.Location,
		&cluster.URL,
		&cluster.Admin.Name,
		&cluster.Admin.Email,
		&cluster.Admin.Password); err != nil {
		return nil, err
	}

	return cluster, nil
}

func (r *ClusterRepository) createClusterTable(ctx context.Context) error {
	defer r.logging(ctx, "CREATE TABLE")()

	query := `CREATE TABLE IF NOT EXISTS cluster (
		id SERIAL NOT NULL PRIMARY KEY,
		location VARCHAR(50) NOT NULL UNIQUE,
		url VARCHAR(50) NOT NULL,
		admin_name VARCHAR(50) NOT NULL,
		admin_email VARCHAR(50) NOT NULL,
		admin_password VARCHAR(100) NOT NULL
	)`

	_, err := r.store.db.Exec(query)

	return err
}

func (r *ClusterRepository) addCluster(ctx context.Context, cluster *models.Cluster) error {
	defer r.logging(ctx, "CREATE cluster for location: "+cluster.Location)()

	if cluster.Admin == nil {
		return errors.New("admin for cluster not set")
	}

	clusterAdmin, err := models.NewAccount(cluster.Admin.Name, cluster.Admin.Email, cluster.Admin.Password)
	if err != nil {
		return err
	}

	if _, err := r.store.Cluster().FindByLocation(ctx, cluster.Location); err == nil {
		r.store.logger.With(
			"table", "cluster",
		).Infof("cluster for location: %v already exists", cluster.Location)

		return nil
	}

	if err := clusterAdmin.Validate(); err != nil {
		return err
	}

	if err := clusterAdmin.BeforeCreate(); err != nil {
		return err
	}

	cluster.Admin = clusterAdmin

	if err := r.store.Cluster().Create(ctx, cluster); err != nil {
		return err
	}

	return nil
}

func (r *ClusterRepository) logging(ctx context.Context, query string) func() {
	sugar := r.store.logger.With(
		"table", "cluster",
	)
	start := time.Now().UTC()
	sugar.Infof("started query %s", query)

	return func() {
		sugar.Infof("complited query %s in %v",
			query,
			time.Since(start))
	}
}
