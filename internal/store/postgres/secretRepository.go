package postgres

import (
	"IAAS/internal/models"
	"IAAS/internal/store"
	"database/sql"
	"time"
)

type SecretRepository struct {
	store *Store
}

func (r *SecretRepository) GetByType(t string) (*models.Secret, error) {
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

func (r *SecretRepository) Init(secret *models.Secret) error {
	if err := r.createSecretTable(); err != nil {
		return err
	}
	if err := r.storeSecret(secret); err != nil {
		return err
	}
	return nil
}

func (r *SecretRepository) createSecretTable() error {
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

func (r *SecretRepository) storeSecret(secret *models.Secret) error {
	query := `
	INSERT INTO secret 
	(value, type, created_at, updated_at)
	SELECT $1, $2, $3, $4
	WHERE
    NOT EXISTS (
        SELECT type FROM secret WHERE type = $5
    );`

	secret.CreatedAt = time.Now()
	secret.UpdatedAt = time.Now()

	_, err := r.store.db.Query(
		query,
		secret.Value,
		secret.Type,
		secret.CreatedAt,
		secret.UpdatedAt,
		secret.Type,
	)
	if err != nil {
		return err
	}

	return nil
}
