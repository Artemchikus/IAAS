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

type apiFunc func(http.ResponseWriter, *http.Request) error

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

func wirteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)

	s.router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin)).Methods("POST")
	s.router.HandleFunc("/account", makeHTTPHandleFunc(s.handleGetAllAccounts)).Methods("GET")
	s.router.HandleFunc("/account", makeHTTPHandleFunc(s.handleCreateAccount)).Methods("POST")
	s.router.HandleFunc("/refreshToken", makeHTTPHandleFunc(s.handleRefreshToken)).Methods("GET")

	private := s.router.PathPrefix("/account/{id}").Subrouter()

	private.Use(s.authenticateAccount)

	private.HandleFunc("", makeHTTPHandleFunc(s.handleGetAccountByID)).Methods("GET")
	private.HandleFunc("", makeHTTPHandleFunc(s.handleDeleteAccount)).Methods("DELETE")
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			resp := ErrorResponse{
				Error: err.Error(),
			}

			wirteJSON(w, http.StatusInternalServerError, resp)
		}
	}
}

func getId(r *http.Request) (int, error) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return id, err
	}
	return id, nil
}
