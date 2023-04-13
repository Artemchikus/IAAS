package api

import (
	"IAAS/internal/config"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"context"
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

	ctx := context.WithValue(context.Background(), models.CtxKeyRequestID, "initial-request")

	store := postgres.New(ctx, db, config)

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
