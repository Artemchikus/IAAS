package api

import (
	"IAAS/internal/models"
	"IAAS/internal/store"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {
	req := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	acc, err := s.store.Account().FindByEmail(r.Context(), req.Email)
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
		return
	}

	if !acc.ValidatePassword(req.Password) {
		s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
		return
	}

	secret, err := s.store.Secret().FindByType(r.Context(), "jwt")
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	jwtToken, err := createJWT(acc, secret.Value)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	refreshToken, err := createRefreshToken(acc, secret.Value)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := s.store.Account().UpdateRefreshToken(r.Context(), acc.RefreshToken, refreshToken); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	resp := &LoginResponse{
		JWT:     jwtToken,
		Refresh: refreshToken,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.store.Account().GetAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	resp := make(GetAllAccountsResponse, 0)

	for _, acc := range accounts {
		r := GetAccountResponse{
			ID:    acc.ID,
			Email: acc.Email,
			Name:  acc.Name,
			Role:  acc.Role,
		}
		resp = append(resp, &r)
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetAccountByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	accId, err := strconv.Atoi(vars["account_id"])
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	if r.Context().Value(models.CtxKeyAccount).(*models.Account).ID != accId && r.Context().Value(models.CtxKeyAccount).(*models.Account).Role != "admin" {
		s.error(w, r, http.StatusUnauthorized, errNotAutheticated)
		return
	}

	account, err := s.store.Account().FindByID(r.Context(), accId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := GetAccountResponse{
		ID:    account.ID,
		Name:  account.Name,
		Email: account.Email,
		Role:  account.Role,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	account := models.NewAccount(req.Name, req.Email, req.Password, "user")

	if _, err := s.store.Account().FindByEmail(r.Context(), account.Email); err == nil {
		s.error(w, r, http.StatusConflict, errEmailAlreadyExists)
		return
	}

	if err := account.Validate(); err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	if err := account.BeforeCreate(); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := s.store.Account().Create(r.Context(), account); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	resp := CreateAccountResponse{
		ID:    account.ID,
		Email: account.Email,
		Name:  account.Name,
		Role:  account.Role,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	accId, err := strconv.Atoi(vars["account_id"])
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	if r.Context().Value(models.CtxKeyAccount).(*models.Account).ID != accId && r.Context().Value(models.CtxKeyAccount).(*models.Account).Role != "admin" {
		s.error(w, r, http.StatusUnauthorized, errNotAutheticated)
		return
	}

	if err := s.store.Account().Delete(r.Context(), accId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteAccountResponse{
		DeletedID: accId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenStr := r.Header.Get("x-refresh-token")

	secret, err := s.store.Secret().FindByType(r.Context(), "jwt")
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	acc, err := getAccFromToken(r.Context(), s.store, refreshTokenStr, secret.Value)
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	if acc.RefreshToken != refreshTokenStr {
		s.error(w, r, http.StatusUnauthorized, errIncorrectToken)
		return
	}

	jwtToken, err := createJWT(acc, secret.Value)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	refreshToken, err := createRefreshToken(acc, secret.Value)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := s.store.Account().UpdateRefreshToken(r.Context(), acc.RefreshToken, refreshToken); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	resp := &RefreshTokenResponse{
		Refresh: refreshToken,
		JWT:     jwtToken,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetAllClusters(w http.ResponseWriter, r *http.Request) {
	clusters, err := s.store.Cluster().GetAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	resp := make(GetAllClustersResponse, 0)

	for _, cl := range clusters {
		r := GetClusterResponse{
			ID:       cl.ID,
			Location: cl.Location,
			URL:      cl.URL,
		}
		resp = append(resp, &r)
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetClusterByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	clusterId, err := strconv.Atoi(vars["cluster_id"])
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	cluster, err := s.store.Cluster().FindByID(r.Context(), clusterId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := GetClusterResponse{
		ID:       cluster.ID,
		Location: cluster.Location,
		URL:      cluster.URL,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleCreateCluster(w http.ResponseWriter, r *http.Request) {
	req := new(CreateClusterRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	cluster := models.NewCluster(req.Admin.Name, req.Admin.Email, req.Admin.Password, req.Admin.DomainID, req.Admin.ProjectID, req.Url, req.Location)

	cluster.Admin.CluserRole = "admin"
	cluster.Admin.ID = req.Admin.ID

	if _, err := s.store.Cluster().FindByLocation(r.Context(), cluster.Location); err == nil {
		s.error(w, r, http.StatusConflict, errLocationAlreadyExists)
		return
	}

	if err := s.store.Cluster().Create(r.Context(), cluster); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := s.store.ClusterUser().Create(r.Context(), cluster.Admin); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := s.fetcher.UpdateClusterMap(r.Context(), s.store); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	resp := CreateClusterResponse{
		ID:        cluster.ID,
		Location:  cluster.Location,
		URL:       cluster.URL,
		AdminName: cluster.Admin.Name,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleDeleteCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	clusterId, err := strconv.Atoi(vars["cluster_id"])
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := s.store.Cluster().Delete(r.Context(), clusterId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	users, err := s.store.ClusterUser().FindByClusterID(r.Context(), clusterId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	for _, user := range users {
		if err := s.store.ClusterUser().Delete(r.Context(), user.ID); err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}
	}

	if err := s.fetcher.UpdateClusterMap(r.Context(), s.store); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	resp := DeleteClusterResponse{
		DeletedID: clusterId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleRegisterAccountInCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	clusterId, err := strconv.Atoi(vars["cluster_id"])
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyClusterID, clusterId))

	token, err := s.fetcher.Token().GetAdmin(r.Context())
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	acc := r.Context().Value(models.CtxKeyAccount).(*models.Account)

	project := models.NewProject(acc.Email, "project for user "+acc.Email)

	if err := s.fetcher.Project().Create(r.Context(), project); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	user := models.NewClusterUser(acc.Name, acc.Email, acc.Password, project.ID, project.DomainID, "common user")

	if err := s.fetcher.User().Create(r.Context(), user); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := s.store.ClusterUser().Create(r.Context(), user); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, user)
}

func (s *server) handleGetFlavorByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	flavorId := vars["flavor_id"]

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	flavor, err := s.fetcher.Flavor().FetchByID(r.Context(), flavorId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, flavor)
}

func (s *server) handleCreateFlavor(w http.ResponseWriter, r *http.Request) {
	req := new(CreateFlavorRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	flavor := models.NewFlavor(req.Name, req.Description, req.VCPUs, req.Ephemeral, req.Disk, req.RAM, req.Swap, req.IsPublic, req.RXTXFactor)

	if err := s.fetcher.Flavor().Create(r.Context(), flavor); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, flavor)
}

func (s *server) handleDeleteFlavor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	flavorId := vars["flavor_id"]

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.Flavor().Delete(r.Context(), flavorId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteOpenstackResurceResponse{
		DeletedID: flavorId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetClusterUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userId := vars["user_id"]

	user, err := s.store.ClusterUser().FindByID(r.Context(), userId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, user)
}

func (s *server) handleGetAllClusterUsers(w http.ResponseWriter, r *http.Request) {
	clusterUsers, err := s.store.ClusterUser().FindByClusterID(r.Context(), r.Context().Value(models.CtxKeyClusterID).(int))
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, clusterUsers)
}

func (s *server) handleCreateClusterUser(w http.ResponseWriter, r *http.Request) {
	req := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	user := models.NewClusterUser(req.Name, req.Email, req.Password, req.ProjectID, req.DomainID, req.Description)

	if err := s.fetcher.User().Create(r.Context(), user); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := s.store.ClusterUser().Create(r.Context(), user); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, user)
}

func (s *server) handleDeleteClusterUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userId := vars["user_id"]

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.User().Delete(r.Context(), userId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	if err := s.store.ClusterUser().Delete(r.Context(), userId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteOpenstackResurceResponse{
		DeletedID: userId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func createJWT(account *models.Account, secret string) (string, error) {
	claims := &jwt.MapClaims{
		"ExpiresAt":    time.Now().UTC().Add(time.Minute * 30).Unix(),
		"AccountEmail": account.Email,
		"AccountRole":  account.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func createRefreshToken(account *models.Account, secret string) (string, error) {
	claims := &jwt.MapClaims{
		"ExpiresAt":    time.Now().UTC().Add(time.Hour * 24).Unix(),
		"AccountEmail": account.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func getAccFromToken(ctx context.Context, s store.Storage, tokenStr, secret string) (*models.Account, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errUnexpectedSigningMethod
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, errIncorrectToken
	}

	if !token.Valid {
		return nil, errNotAutheticated
	}

	claims := token.Claims.(jwt.MapClaims)
	if float64(time.Now().UTC().Unix()) >= claims["ExpiresAt"].(float64) {
		return nil, errTokenExpired
	}

	email := claims["AccountEmail"].(string)

	acc, err := s.Account().FindByEmail(ctx, email)
	if err != nil {
		return nil, errIncorrectToken
	}

	return acc, nil
}
