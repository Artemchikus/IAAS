package postgres

import (
	"IAAS/internal/models"
	"IAAS/internal/store"
	"database/sql"
	"time"
)

type AccountRepository struct {
	store *Store
}

func (r *AccountRepository) Create(account *models.Account) error {
	if err := account.Validate(); err != nil {
		return err
	}

	if err := account.BeforeCreate(); err != nil {
		return err
	}

	query := `
	INSERT INTO account 
	(name, email, encrypted_password, created_at, updated_at, refresh_token)  
	values ($1, $2, $3, $4, $5, $6) RETURNING *`

	row, err := r.store.db.Query(
		query,
		account.Name,
		account.Email,
		account.EncryptedPassword,
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

func (r *AccountRepository) Delete(id int) error {
	row, err := r.store.db.Query("DELETE FROM account WHERE id = $1", id)
	if err != nil {
		return store.ErrRecordNotFound
	}
	defer row.Close()

	return nil
}

func (r *AccountRepository) Update(a *models.Account) error {
	return nil
}

func (r *AccountRepository) FindByEmail(email string) (*models.Account, error) {
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

func (r *AccountRepository) FindByID(id int) (*models.Account, error) {
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

func (r *AccountRepository) Init() error {
	return r.createAccountTable()
}

func (r *AccountRepository) GetAll() ([]*models.Account, error) {
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
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.RefreshToken); err != nil {
		return nil, err
	}

	return account, nil
}

func (r *AccountRepository) UpdateRefreshToken(old, new string, time time.Time) error {
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

func (r *AccountRepository) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id SERIAL NOT NULL PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		email VARCHAR(50) NOT NULL UNIQUE,
		encrypted_password VARCHAR(100) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		refresh_token VARCHAR(255)
	)`

	_, err := r.store.db.Exec(query)

	return err
}
