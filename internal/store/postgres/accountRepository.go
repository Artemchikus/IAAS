package postgres

import (
	"IAAS/internal/models"
	"IAAS/internal/store"
	"context"
	"database/sql"
	"log"
	"time"
)

type AccountRepository struct {
	store *Store
}

func (r *AccountRepository) Create(ctx context.Context, account *models.Account) error {
	defer r.logging(ctx, "CREATE")()

	query := `
	INSERT INTO account 
	(name, email, encrypted_password, role, created_at, updated_at, refresh_token)  
	values ($1, $2, $3, $4, $5, $6, $7) RETURNING *`

	row, err := r.store.db.Query(
		query,
		account.Name,
		account.Email,
		account.EncryptedPassword,
		account.Role,
		account.CreatedAt,
		account.UpdatedAt,
		account.RefreshToken)
	if err != nil {
		return err
	}
	defer row.Close()

	for row.Next() {
		acc, err := scanIntoAccount(row)
		if err != nil {
			return err
		}
		account.ID = acc.ID
	}

	return nil
}

func (r *AccountRepository) Delete(ctx context.Context, id int) error {
	defer r.logging(ctx, "DELETE")()

	row, err := r.store.db.Query("DELETE FROM account WHERE id = $1", id)
	if err != nil {
		return store.ErrRecordNotFound
	}
	defer row.Close()

	return nil
}

func (r *AccountRepository) Update(ctx context.Context, a *models.Account) error {
	return nil
}

func (r *AccountRepository) FindByEmail(ctx context.Context, email string) (*models.Account, error) {
	defer r.logging(ctx, "FIND BY email")()

	rows, err := r.store.db.Query("SELECT * FROM account WHERE email = $1", email)
	if err != nil {
		return nil, store.ErrRecordNotFound
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, store.ErrRecordNotFound
}

func (r *AccountRepository) FindByID(ctx context.Context, id int) (*models.Account, error) {
	defer r.logging(ctx, "FIND BY id")()

	rows, err := r.store.db.Query("SELECT * FROM account WHERE id = $1", id)
	if err != nil {
		return nil, store.ErrRecordNotFound
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, store.ErrRecordNotFound
}

func (r *AccountRepository) Init(ctx context.Context, admin *models.Account) error {
	defer r.logging(ctx, "INIT")()

	if err := r.createAccountTable(ctx); err != nil {
		return err
	}

	if err := r.createAdmin(ctx, admin); err != nil {
		return err
	}

	return nil
}

func (r *AccountRepository) GetAll(ctx context.Context) ([]*models.Account, error) {
	defer r.logging(ctx, "GET ALL")()

	rows, err := r.store.db.Query("SELECT * FROM account")
	if err != nil {
		return nil, store.ErrRecordNotFound
	}
	defer rows.Close()

	accounts := []*models.Account{}

	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, store.ErrRecordNotFound
		}
		accounts = append(accounts, account)
	}

	if len(accounts) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*models.Account, error) {
	account := new(models.Account)
	if err := rows.Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.EncryptedPassword,
		&account.Role,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.RefreshToken); err != nil {
		return nil, err
	}

	return account, nil
}

func (r *AccountRepository) UpdateRefreshToken(ctx context.Context, old, new string, time time.Time) error {
	defer r.logging(ctx, "UPDATE refresh_token")()

	query := `UPDATE account
	SET refresh_token = $1,
		updated_at = $2
	WHERE refresh_token = $3;`

	row, err := r.store.db.Query(
		query,
		new,
		time,
		old)
	if err != nil {
		return err
	}
	defer row.Close()

	return nil
}

func (r *AccountRepository) createAccountTable(ctx context.Context) error {
	defer r.logging(ctx, "CREATE TABLE")()

	query := `CREATE TABLE IF NOT EXISTS account (
		id SERIAL NOT NULL PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		email VARCHAR(50) NOT NULL UNIQUE,
		encrypted_password VARCHAR(100) NOT NULL,
		role VARCHAR(50) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		refresh_token VARCHAR(255)
	)`

	_, err := r.store.db.Exec(query)

	return err
}

func (r *AccountRepository) createAdmin(ctx context.Context, admin *models.Account) error {
	defer r.logging(ctx, "CREATE admin")()

	adm, err := models.NewAccount(admin.Name, admin.Email, admin.Password)
	if err != nil {
		return err
	}

	log.Println(adm.Email)

	if _, err := r.store.Account().FindByEmail(ctx, adm.Email); err == nil {
		r.store.logger.With(
			"table", "account",
		).Info("initial admin already exists")

		return nil
	}

	if err := admin.Validate(); err != nil {
		return err
	}

	if err := admin.BeforeCreate(); err != nil {
		return err
	}

	admin.Role = "admin"

	if err := r.store.Account().Create(ctx, admin); err != nil {
		return err
	}

	return nil
}

func (r *AccountRepository) logging(ctx context.Context, query string) func() {
	sugar := r.store.logger.With(
		"table", "account",
	)
	start := time.Now()
	sugar.Infof("started query %s", query)

	return func() {
		sugar.Infof("complited query %s in %v",
			query,
			time.Since(start))
	}
}
