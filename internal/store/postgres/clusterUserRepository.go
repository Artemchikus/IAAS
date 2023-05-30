package postgres

import (
	"IAAS/internal/models"
	"IAAS/internal/store"
	"context"
	"database/sql"
	"time"
)

type ClusterUserRepository struct {
	store *Store
}

func (r *ClusterUserRepository) Create(ctx context.Context, clusterUser *models.ClusterUser) error {
	defer r.logging(ctx, "CREATE")()

	query := `
	INSERT INTO clusterUser 
	(id, name, email, password, role, domain_id, account_id, cluster_id, project_id, description)  
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *`

	row, err := r.store.db.Query(
		query,
		clusterUser.ID,
		clusterUser.Name,
		clusterUser.Email,
		clusterUser.Password,
		clusterUser.CluserRole,
		clusterUser.DomainID,
		clusterUser.AccountID,
		clusterUser.ClusterID,
		clusterUser.ProjectID,
		clusterUser.Description)
	if err != nil {
		return err
	}
	defer row.Close()

	for row.Next() {
		acc, err := scanIntoClusterUser(row)
		if err != nil {
			return err
		}
		clusterUser.ID = acc.ID
	}

	return nil
}

func (r *ClusterUserRepository) Delete(ctx context.Context, id string) error {
	defer r.logging(ctx, "DELETE")()

	row, err := r.store.db.Query("DELETE FROM clusterUser WHERE id = $1", id)
	if err != nil {
		return store.ErrRecordNotFound
	}
	defer row.Close()

	return nil
}

func (r *ClusterUserRepository) Update(ctx context.Context, a *models.ClusterUser) error {
	return nil
}

func (r *ClusterUserRepository) FindByAccountID(ctx context.Context, accId int) ([]*models.ClusterUser, error) {
	defer r.logging(ctx, "FIND BY account_id")()

	rows, err := r.store.db.Query("SELECT * FROM clusterUser WHERE account_id = $1", accId)
	if err != nil {
		return nil, store.ErrRecordNotFound
	}
	defer rows.Close()

	clusterUsers := []*models.ClusterUser{}

	for rows.Next() {
		clusterUser, err := scanIntoClusterUser(rows)
		if err != nil {
			return nil, store.ErrRecordNotFound
		}
		clusterUsers = append(clusterUsers, clusterUser)
	}

	if len(clusterUsers) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return clusterUsers, nil
}

func (r *ClusterUserRepository) FindByClusterID(ctx context.Context, clusterId int) ([]*models.ClusterUser, error) {
	defer r.logging(ctx, "FIND BY cluster_id")()

	rows, err := r.store.db.Query("SELECT * FROM clusterUser WHERE cluster_id = $1", clusterId)
	if err != nil {
		return nil, store.ErrRecordNotFound
	}
	defer rows.Close()

	clusterUsers := []*models.ClusterUser{}

	for rows.Next() {
		clusterUser, err := scanIntoClusterUser(rows)
		if err != nil {
			return nil, store.ErrRecordNotFound
		}
		clusterUsers = append(clusterUsers, clusterUser)
	}

	if len(clusterUsers) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return clusterUsers, nil
}

func (r *ClusterUserRepository) FindByID(ctx context.Context, id string) (*models.ClusterUser, error) {
	defer r.logging(ctx, "FIND BY id")()

	rows, err := r.store.db.Query("SELECT * FROM clusterUser WHERE id = $1", id)
	if err != nil {
		return nil, store.ErrRecordNotFound
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoClusterUser(rows)
	}

	return nil, store.ErrRecordNotFound
}

func (r *ClusterUserRepository) FindByEmailAndClusterID(ctx context.Context, email string, clusterId int) (*models.ClusterUser, error) {
	defer r.logging(ctx, "FIND BY email AND cluster_id")()

	rows, err := r.store.db.Query("SELECT * FROM clusterUser WHERE email = $1 AND cluster_id = $2", email, clusterId)
	if err != nil {
		return nil, store.ErrRecordNotFound
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoClusterUser(rows)
	}

	return nil, store.ErrRecordNotFound
}

func (r *ClusterUserRepository) Init(ctx context.Context, clusters []*models.Cluster) error {
	defer r.logging(ctx, "INIT")()

	if err := r.createClusterUserTable(ctx); err != nil {
		return err
	}

	for _, cluster := range clusters {
		cluster.Admin.CluserRole = "admin"

		if _, err := r.store.ClusterUser().FindByEmailAndClusterID(ctx, cluster.Admin.Email, cluster.ID); err == nil {
			r.store.logger.With(
				"table", "clusterUser",
			).Infof("clusterUser for email: %v and clusteId: %v ralready exists", cluster.Admin.Email, cluster.Admin.ClusterID)

			continue
		}

		if err := r.Create(ctx, cluster.Admin); err != nil {
			return err
		}
	}

	return nil
}

func (r *ClusterUserRepository) GetAll(ctx context.Context) ([]*models.ClusterUser, error) {
	defer r.logging(ctx, "GET ALL")()

	rows, err := r.store.db.Query("SELECT * FROM clusterUser")
	if err != nil {
		return nil, store.ErrRecordNotFound
	}
	defer rows.Close()

	clusterUsers := []*models.ClusterUser{}

	for rows.Next() {
		clusterUser, err := scanIntoClusterUser(rows)
		if err != nil {
			return nil, store.ErrRecordNotFound
		}
		clusterUsers = append(clusterUsers, clusterUser)
	}

	if len(clusterUsers) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return clusterUsers, nil
}

func scanIntoClusterUser(rows *sql.Rows) (*models.ClusterUser, error) {
	clusterUser := new(models.ClusterUser)
	if err := rows.Scan(
		&clusterUser.ID,
		&clusterUser.Name,
		&clusterUser.Email,
		&clusterUser.Password,
		&clusterUser.CluserRole,
		&clusterUser.DomainID,
		&clusterUser.AccountID,
		&clusterUser.ClusterID,
		&clusterUser.ProjectID,
		&clusterUser.Description); err != nil {
		return nil, err
	}

	return clusterUser, nil
}

func (r *ClusterUserRepository) createClusterUserTable(ctx context.Context) error {
	defer r.logging(ctx, "CREATE TABLE")()

	query := `CREATE TABLE IF NOT EXISTS clusterUser (
		id VARCHAR(50) NOT NULL PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		email VARCHAR(50) NOT NULL,
		password VARCHAR(100) NOT NULL,
		role VARCHAR(50) NOT NULL,
		domain_id VARCHAR(50) NOT NULL,
		account_id INTEGER,
		cluster_id INTEGER NOT NULL,
		project_id VARCHAR(50) NOT NULL,
		description VARCHAR(50)
	)`

	_, err := r.store.db.Exec(query)

	return err
}

func (r *ClusterUserRepository) logging(ctx context.Context, query string) func() {
	sugar := r.store.logger.With(
		"table", "clusterUser",
	)
	start := time.Now().UTC()
	sugar.Infof("started query %s", query)

	return func() {
		sugar.Infof("complited query %s in %v",
			query,
			time.Since(start))
	}
}
