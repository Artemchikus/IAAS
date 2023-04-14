package postgres

import (
	"IAAS/internal/models"
	"IAAS/internal/store"
	"context"
	"database/sql"
	"time"
)

type SecretRepository struct {
	store *Store
}

func (r *SecretRepository) FindByType(ctx context.Context, t string) (*models.Secret, error) {
	defer r.logging(ctx, "GET BY type")()

	rows, err := r.store.db.Query("SELECT * FROM secret WHERE type = $1", t)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoSecret(rows)
	}

	return nil, store.ErrRecordNotFound
}

func (r *SecretRepository) Init(ctx context.Context, secret *models.Secret) error {
	defer r.logging(ctx, "INIT")()

	if err := r.createSecretTable(ctx); err != nil {
		return err
	}
	if err := r.storeSecret(ctx, secret); err != nil {
		return err
	}
	return nil
}

func (r *SecretRepository) createSecretTable(ctx context.Context) error {
	defer r.logging(ctx, "CREATE TABLE")()

	query := `CREATE TABLE IF NOT EXISTS secret (
		id SERIAL NOT NULL PRIMARY KEY,
		type VARCHAR(50) NOT NULL UNIQUE,
		value VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	)`

	_, err := r.store.db.Exec(query)

	return err
}

func scanIntoSecret(rows *sql.Rows) (*models.Secret, error) {
	secret := new(models.Secret)
	if err := rows.Scan(
		&secret.ID,
		&secret.Type,
		&secret.Value,
		&secret.CreatedAt,
		&secret.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return secret, nil
}

func (r *SecretRepository) storeSecret(ctx context.Context, secret *models.Secret) error {
	defer r.logging(ctx, "INSERT secret")()

	_, err := r.store.secretRepository.FindByType(ctx, secret.Type)
	if err == nil {
		r.store.logger.With(
			"table", "secret",
			"request_id", ctx.Value(models.CtxKeyRequestID),
		).Info("initital secret already exists")
		return nil
	}

	query := `
	INSERT INTO secret 
	(value, type, created_at, updated_at)
	values ($1, $2, $3, $4)`

	secret.CreatedAt = time.Now()
	secret.UpdatedAt = time.Now()

	_, err = r.store.db.Query(
		query,
		secret.Value,
		secret.Type,
		secret.CreatedAt,
		secret.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *SecretRepository) logging(ctx context.Context, query string) func() {
	sugar := r.store.logger.With(
		"table", "secret",
		"request_id", ctx.Value(models.CtxKeyRequestID),
	)
	start := time.Now()
	sugar.Infof("started query %s", query)

	return func() {
		sugar.Infof("complited query %s in %v",
			query,
			time.Since(start))
	}
}
