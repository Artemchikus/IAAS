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

	account := models.NewAccount(req.Name, req.Email, req.Password, "member")

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

	acc := r.Context().Value(models.CtxKeyAccount).(*models.Account)

	if _, err := s.store.ClusterUser().FindByEmailAndClusterID(r.Context(), acc.Email, clusterId); err == nil {
		s.error(w, r, http.StatusUnprocessableEntity, errAlreadyInCluster)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyClusterID, clusterId))

	token, err := s.fetcher.Token().GetAdmin(r.Context())
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	project := models.NewProject(acc.Email, "project for user "+acc.Email)

	if err := s.fetcher.Project().Create(r.Context(), project); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	user := models.NewClusterUser(acc.Name, acc.Email, acc.EncryptedPassword, project.ID, project.DomainID, "common user")
	user.AccountID = acc.ID
	user.CluserRole = acc.Role

	if err := s.fetcher.User().Create(r.Context(), user); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	role, err := s.fetcher.Role().FetchByName(r.Context(), acc.Role)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := s.fetcher.User().AssignRoleToProject(r.Context(), user.ProjectID, role.ID, user.ID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := s.store.ClusterUser().Create(r.Context(), user); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, user)
}

func (s *server) handleGetFloatingIps(w http.ResponseWriter, r *http.Request) {
	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	floatingIps, err := s.fetcher.FloatingIp().FetchAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, floatingIps)
}

func (s *server) handleGetFloatingIpByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	floatingId, ok := vars["floatingip_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	floatingIp, err := s.fetcher.FloatingIp().FetchByID(r.Context(), floatingId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, floatingIp)
}

func (s *server) handleCreateFloatingIp(w http.ResponseWriter, r *http.Request) {
	req := new(CreateFloatingIpRequest)
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

	floatingIp := models.NewFloatingIP(req.NetworkID, req.Description)

	if err := s.fetcher.FloatingIp().Create(r.Context(), floatingIp); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, floatingIp)
}

func (s *server) handleDeleteFloatingIp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	floatingId, ok := vars["floatingip_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.FloatingIp().Delete(r.Context(), floatingId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteOpenstackResurceResponse{
		DeletedID: floatingId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleAddIPToPort(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	floatingId, ok := vars["floatingip_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	req := new(AddFloatingIPRequest)
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

	if err := s.fetcher.FloatingIp().AddToPort(r.Context(), floatingId, req.PortID); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := AddFloatingIPResponse{
		FloatingIpID: floatingId,
		PortID:       req.PortID,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetImages(w http.ResponseWriter, r *http.Request) {
	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	images, err := s.fetcher.Image().FetchAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, images)
}

func (s *server) handleGetImageByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	imageId, ok := vars["image_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	image, err := s.fetcher.Image().FetchByID(r.Context(), imageId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, image)
}

func (s *server) handleCreateImage(w http.ResponseWriter, r *http.Request) {
	req := new(CreateImageRequest)
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

	image := models.NewImage(req.Name, req.DiskFormat, req.ContainerFormat, req.Visibility)

	if err := s.fetcher.Image().Create(r.Context(), image); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, image)
}

func (s *server) handleDeleteImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	imageId, ok := vars["image_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.Image().Delete(r.Context(), imageId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteOpenstackResurceResponse{
		DeletedID: imageId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleUploadImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	imageId, ok := vars["image_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	var data []byte

	n, err := r.Body.Read(data)
	if err != nil || n == 0 {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.Image().Upload(r.Context(), data, imageId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := UploadImageResponse{
		ImageID: imageId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetVolumes(w http.ResponseWriter, r *http.Request) {
	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	volumes, err := s.fetcher.Volume().FetchAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, volumes)
}

func (s *server) handleGetVolumeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	volumeId, ok := vars["volume_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	volume, err := s.fetcher.Volume().FetchByID(r.Context(), volumeId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, volume)
}

func (s *server) handleCreateVolume(w http.ResponseWriter, r *http.Request) {
	req := new(CreateVolumeRequest)
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

	volume := models.NewVolume(req.Description, req.Name, req.Bootable, req.Size)

	if err := s.fetcher.Volume().Create(r.Context(), volume); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, volume)
}

func (s *server) handleDeleteVolume(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	volumeId, ok := vars["volume_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.Volume().Delete(r.Context(), volumeId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteOpenstackResurceResponse{
		DeletedID: volumeId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetNetworks(w http.ResponseWriter, r *http.Request) {
	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	networks, err := s.fetcher.Network().FetchAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, networks)
}

func (s *server) handleGetNetworkByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	networkId, ok := vars["network_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	network, err := s.fetcher.Network().FetchByID(r.Context(), networkId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, network)
}

func (s *server) handleCreatePublicNetwork(w http.ResponseWriter, r *http.Request) {
	req := new(CreateNetworkRequest)
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

	network := models.NewNetwork(req.Description, req.Name, req.ProjectID, req.MTU, req.External)

	if err := s.fetcher.Network().Create(r.Context(), network); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, network)
}

func (s *server) handleCreatePrivateNetwork(w http.ResponseWriter, r *http.Request) {
	req := new(CreateNetworkRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	token, err := s.fetcher.Token().GetAdmin(r.Context())
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	acc := r.Context().Value(models.CtxKeyAccount).(*models.Account)
	clusterId := r.Context().Value(models.CtxKeyClusterID).(int)

	user, err := s.store.ClusterUser().FindByEmailAndClusterID(r.Context(), acc.Email, clusterId)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	network := models.NewNetwork(req.Description, req.Name, user.ProjectID, req.MTU, req.External)

	if err := s.fetcher.Network().Create(r.Context(), network); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, network)
}

func (s *server) handleDeletePrivateNetwork(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	networkId, ok := vars["network_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.Network().Delete(r.Context(), networkId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteOpenstackResurceResponse{
		DeletedID: networkId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleDeletePublicNetwork(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	networkId, ok := vars["network_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.Network().Delete(r.Context(), networkId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteOpenstackResurceResponse{
		DeletedID: networkId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetServers(w http.ResponseWriter, r *http.Request) {
	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	servers, err := s.fetcher.Server().FetchAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, servers)
}

func (s *server) handleGetServerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	serverId, ok := vars["server_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	server, err := s.fetcher.Server().FetchByID(r.Context(), serverId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, server)
}

func (s *server) handleCreateServer(w http.ResponseWriter, r *http.Request) {
	req := new(CreateServerRequest)
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

	server := models.NewServer(req.ImageID, req.KeyID, req.Name, req.SecurityGroupID, req.PrivateNetworkID, req.FlavorID)

	if err := s.fetcher.Server().Create(r.Context(), server); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, server)
}

func (s *server) handleDeleteServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	serverId, ok := vars["server_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.Server().Delete(r.Context(), serverId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteOpenstackResurceResponse{
		DeletedID: serverId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleStopServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	serverId, ok := vars["server_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.Server().Stop(r.Context(), serverId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := StopServerResponse{
		StopedID: serverId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleStartServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	serverId, ok := vars["server_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.Server().Start(r.Context(), serverId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := StartServerResponse{
		StartedID: serverId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleAttachVolToServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	serverId, ok := vars["server_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	req := new(AttachVolumeRequest)
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

	if err := s.fetcher.Server().AttachVolume(r.Context(), serverId, req.VolumeID); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := AttachVolumeResponse{
		ServerID: serverId,
		VolumeID: req.VolumeID,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetRoles(w http.ResponseWriter, r *http.Request) {
	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	roles, err := s.fetcher.Role().FetchAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, roles)
}

func (s *server) handleGetRoleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	roleId, ok := vars["role_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	role, err := s.fetcher.Role().FetchByID(r.Context(), roleId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, role)
}

func (s *server) handleCreateRole(w http.ResponseWriter, r *http.Request) {
	req := new(CreateRoleRequest)
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

	role := models.NewRole(req.Description, req.Name)

	if err := s.fetcher.Role().Create(r.Context(), role); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, role)
}

func (s *server) handleDeleteRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	roleId, ok := vars["role_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.Role().Delete(r.Context(), roleId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteOpenstackResurceResponse{
		DeletedID: roleId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetSecurityGroups(w http.ResponseWriter, r *http.Request) {
	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	secGroups, err := s.fetcher.SecurityGroup().FetchAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, secGroups)
}

func (s *server) handleGetSecurityGroupByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	secGroupId, ok := vars["securitygroup_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	secGroup, err := s.fetcher.SecurityGroup().FetchByID(r.Context(), secGroupId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, secGroup)
}

func (s *server) handleCreateSecurityGroup(w http.ResponseWriter, r *http.Request) {
	req := new(CreateSecurityGroupRequest)
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

	secGroup := models.NewSecurityGroup(req.Description, req.Name)

	if err := s.fetcher.SecurityGroup().Create(r.Context(), secGroup); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, secGroup)
}

func (s *server) handleDeleteSecurityGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	secGroupId, ok := vars["securitygroup_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.SecurityGroup().Delete(r.Context(), secGroupId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteOpenstackResurceResponse{
		DeletedID: secGroupId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetSubnets(w http.ResponseWriter, r *http.Request) {
	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	subnets, err := s.fetcher.Subnet().FetchAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, subnets)
}

func (s *server) handleGetSubnetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	subnetId, ok := vars["subnet_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	subnet, err := s.fetcher.Subnet().FetchByID(r.Context(), subnetId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, subnet)
}

func (s *server) handleCreateSubnet(w http.ResponseWriter, r *http.Request) {
	req := new(CreateSubnetRequest)
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

	subnet := models.NewSubnet(req.CIDR, req.Description, req.Name, req.NetworkID, req.GatewayIp, req.EnableDHCP, req.AllocationPools)

	if err := s.fetcher.Subnet().Create(r.Context(), subnet); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, subnet)
}

func (s *server) handleDeleteSubnet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	subnetId, ok := vars["subnet_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.Subnet().Delete(r.Context(), subnetId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteOpenstackResurceResponse{
		DeletedID: subnetId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetSecurityRules(w http.ResponseWriter, r *http.Request) {
	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	secRules, err := s.fetcher.SecurityRule().FetchAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, secRules)
}

func (s *server) handleGetSecurityRuleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	secRuleId, ok := vars["securityrule_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	secRule, err := s.fetcher.SecurityRule().FetchByID(r.Context(), secRuleId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, secRule)
}

func (s *server) handleCreateSecurityRule(w http.ResponseWriter, r *http.Request) {
	req := new(CreateSecurityRuleRequest)
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

	secRule := models.NewSecurityRule(req.Ethertype, req.Direction, req.Protocol, req.RemoteIpPrefix, req.Description, req.SecurityGroupID, req.PortRangeMax, req.PortRangeMin)

	if err := s.fetcher.SecurityRule().Create(r.Context(), secRule); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, secRule)
}

func (s *server) handleDeleteSecurityRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	secRuleId, ok := vars["securityrule_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.SecurityRule().Delete(r.Context(), secRuleId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteOpenstackResurceResponse{
		DeletedID: secRuleId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetKeyPairs(w http.ResponseWriter, r *http.Request) {
	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	keyPairs, err := s.fetcher.KeyPair().FetchAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, keyPairs)
}

func (s *server) handleGetKeyPairByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	keyPairId, ok := vars["keypair_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	keyPair, err := s.fetcher.KeyPair().FetchByID(r.Context(), keyPairId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, keyPair)
}

func (s *server) handleCreateKeyPair(w http.ResponseWriter, r *http.Request) {
	req := new(CreateKeyPairRequest)
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

	keyPair := models.NewKeyPair(req.PublicKey, req.Name, req.Type)

	if err := s.fetcher.KeyPair().Create(r.Context(), keyPair); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, keyPair)
}

func (s *server) handleDeleteKeyPair(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	keyPairId, ok := vars["keypair_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.KeyPair().Delete(r.Context(), keyPairId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteOpenstackResurceResponse{
		DeletedID: keyPairId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetProjects(w http.ResponseWriter, r *http.Request) {
	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	projects, err := s.fetcher.Project().FetchAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, projects)
}

func (s *server) handleGetProjectByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	projectId, ok := vars["project_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	project, err := s.fetcher.Project().FetchByID(r.Context(), projectId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, project)
}

func (s *server) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	req := new(CreateProjectRequest)
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

	project := models.NewProject(req.Name, req.Description)

	if err := s.fetcher.Project().Create(r.Context(), project); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, project)
}

func (s *server) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	projectId, ok := vars["project_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.Project().Delete(r.Context(), projectId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteOpenstackResurceResponse{
		DeletedID: projectId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetRouters(w http.ResponseWriter, r *http.Request) {
	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	routers, err := s.fetcher.Router().FetchAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, routers)
}

func (s *server) handleGetRouterByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	routerId, ok := vars["router_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	router, err := s.fetcher.Router().FetchByID(r.Context(), routerId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, router)
}

func (s *server) handleCreateRouter(w http.ResponseWriter, r *http.Request) {
	req := new(CreateRouterRequest)
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

	router := models.NewRouter(req.Description, req.Name, req.ExternalGatewayInfo.NetworkID)

	if err := s.fetcher.Router().Create(r.Context(), router); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, router)
}

func (s *server) handleDeleteRouter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	routerId, ok := vars["router_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	if err := s.fetcher.Router().Delete(r.Context(), routerId); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := DeleteOpenstackResurceResponse{
		DeletedID: routerId,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleAddSubnetToRouter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	routerId, ok := vars["router_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	req := new(RemoveOrAddSubnetRequest)
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

	if err := s.fetcher.Router().AddSubnet(r.Context(), routerId, req.SubnetID); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := RemoveOrAddSubnetResponse{
		RouterID: routerId,
		SubnetID: req.SubnetID,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleRemoveRouterSubnet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	routerId, ok := vars["router_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	req := new(RemoveOrAddSubnetRequest)
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

	if err := s.fetcher.Router().RemoveSubnet(r.Context(), routerId, req.SubnetID); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	resp := RemoveOrAddSubnetResponse{
		RouterID: routerId,
		SubnetID: req.SubnetID,
	}

	s.respond(w, r, http.StatusOK, resp)
}

func (s *server) handleGetPorts(w http.ResponseWriter, r *http.Request) {
	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	ports, err := s.fetcher.Port().FetchAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, ports)
}

func (s *server) handleGetPortByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	portId, ok := vars["port_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	port, err := s.fetcher.Port().FetchByID(r.Context(), portId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, port)
}

func (s *server) handleGetPortsByNetwork(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	networkId, ok := vars["network_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	ports, err := s.fetcher.Port().FetchByNetworkID(r.Context(), networkId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, ports)
}

func (s *server) handleGetPortsByDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deviceId, ok := vars["device_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	ports, err := s.fetcher.Port().FetchByDeviceID(r.Context(), deviceId)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, ports)
}

func (s *server) handleGetFlavors(w http.ResponseWriter, r *http.Request) {
	token, err := s.fetcher.Token().Get(r.Context(), r.Context().Value(models.CtxClusterUser).(*models.ClusterUser))
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), models.CtxKeyToken, token))

	flavors, err := s.fetcher.Flavor().FetchAll(r.Context())
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, flavors)
}

func (s *server) handleGetFlavorByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	flavorId, ok := vars["flavor_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

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

	flavorId, ok := vars["flavor_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

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

	userId, ok := vars["user_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

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

	userId, ok := vars["user_id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, errBadRequest)
		return
	}

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
