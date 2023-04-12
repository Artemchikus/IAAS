package api

import (
	"IAAS/internal/store"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ctxKey int8

const (
	ctxKeyRequestID ctxKey = iota
)

type server struct {
	router *mux.Router
	store  store.Storage
	logger *zap.SugaredLogger
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newServer(store store.Storage) *server {
	zapLog, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer zapLog.Sync()
	sugar := zapLog.Sugar()

	s := &server{
		logger: sugar,
		router: mux.NewRouter(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	resp := ErrorResponse{
		Error: err.Error(),
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
	s.router.HandleFunc("/account", s.handleGetAllAccounts).Methods("GET")
	s.router.HandleFunc("/account", s.handleCreateAccount).Methods("POST")
	s.router.HandleFunc("/refreshToken", s.handleRefreshToken).Methods("GET")

	private := s.router.PathPrefix("/account/{id}").Subrouter()

	private.Use(s.authenticateAccount)

	private.HandleFunc("", s.handleGetAccountByID).Methods("GET")
	private.HandleFunc("", s.handleDeleteAccount).Methods("DELETE")
}

func getId(r *http.Request) (int, error) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return id, err
	}
	return id, nil
}
