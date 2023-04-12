package store

type Storage interface {
	Account() AccountRepository
	Secret() SecretRepository
}
