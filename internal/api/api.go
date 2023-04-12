package api

import (
	"IAAS/internal/config"
	"IAAS/internal/store/postgres"
	"database/sql"
	"log"
	"net/http"
)

func Start(config *config.ApiConfig) error {
	URL := config.DatabaseURL
	db, err := newDB(URL)
	if err != nil {
		return err
	}

	defer db.Close()

	store := postgres.New(db, config)

	srv := newServer(store)

	addr := config.BindAddr

	log.Println("Server is listening on port", addr)

	return http.ListenAndServe(addr, srv)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connection to db is successfull")
	return db, nil
}
