package api

import (
	"IAAS/internal/business"
	"IAAS/internal/store"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type server struct {
	router  *mux.Router
	store   store.Storage
	logger  *zap.SugaredLogger
	fetcher business.Fetcher
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newServer(store store.Storage, fetcher business.Fetcher) *server {
	zapLog, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer zapLog.Sync()
	sugar := zapLog.Sugar()

	s := &server{
		logger:  sugar,
		router:  mux.NewRouter(),
		store:   store,
		fetcher: fetcher,
	}

	s.configureRouter()

	return s
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.logger.Errorf("error: %v", err)

	incErr := incapsulateError(code, err)

	resp := ErrorResponse{
		Error: incErr.Error(),
	}
	s.respond(w, r, code, resp)
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)

	s.router.HandleFunc("/login", s.handleLogin).Methods("POST")
	s.router.HandleFunc("/account", s.handleCreateAccount).Methods("POST")
	s.router.HandleFunc("/token/refresh", s.handleRefreshToken).Methods("GET")

	private := s.router.PathPrefix("/account/{account_id}").Subrouter()

	private.Use(s.authenticateAccount)

	private.HandleFunc("", s.handleGetAccountByID).Methods("GET")
	private.HandleFunc("", s.handleDeleteAccount).Methods("DELETE")
	private.HandleFunc("/cluster", s.handleGetAllClusters).Methods("GET")
	private.HandleFunc("/cluster/{cluster_id}", s.handleGetClusterByID).Methods("GET")

	admin := private.PathPrefix("/admin").Subrouter()

	admin.Use(s.isAdmin)

	admin.HandleFunc("/account", s.handleGetAllAccounts).Methods("GET")
	admin.HandleFunc("/cluster", s.handleCreateCluster).Methods("POST")
	admin.HandleFunc("/cluster/{cluster_id}", s.handleDeleteCluster).Methods("DELETE")

}

func getVars(r *http.Request) (map[string]int, error) {
	vars := mux.Vars(r)

	intVars := make(map[string]int)

	for k, v := range vars {
		intV, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		intVars[k] = intV
	}
	return intVars, nil
}

func incapsulateError(code int, err error) error {
	switch code {
	case http.StatusUnauthorized:
		return err

	case http.StatusBadRequest:
		return errBadRequest

	case http.StatusInternalServerError:
		if err == store.ErrRecordNotFound {
			return err
		}
		return errInternalServerError
	case http.StatusNotFound:
		return err

	case http.StatusConflict:
		return errEmailAlreadyExists

	case http.StatusUnprocessableEntity:
		return err

	}
	return err
}
