package api

import (
	"IAAS/internal/business"
	"IAAS/internal/store"
	"encoding/json"
	"log"
	"net/http"

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

	private := s.router.PathPrefix("").Subrouter()

	private.Use(s.authenticateAccount)

	private.HandleFunc("/account/{account_id}", s.handleGetAccountByID).Methods("GET")
	private.HandleFunc("/account/{account_id}", s.handleDeleteAccount).Methods("DELETE")
	private.HandleFunc("/cluster", s.handleGetAllClusters).Methods("GET")
	private.HandleFunc("/cluster/{cluster_id}", s.handleGetClusterByID).Methods("GET")
	private.HandleFunc("/cluster/{cluster_id}/register", s.handleRegisterAccountInCluster).Methods("HEAD")

	cluster := private.PathPrefix("/cluster/{cluster_id}").Subrouter()

	cluster.Use(s.authenticateClusterUser)

	cluster.HandleFunc("/flavor", s.handleGetFlavors).Methods("GET")
	cluster.HandleFunc("/flavor/{flavor_id}", s.handleGetFlavorByID).Methods("GET")
	cluster.HandleFunc("/floatingIp", s.handleGetFloatingIps).Methods("GET")
	cluster.HandleFunc("/floatingIp/{floatingIp_id}", s.handleGetFloatingIpByID).Methods("GET")
	cluster.HandleFunc("/floatingIp/{floatingIp_id}", s.handleAddIPToPort).Methods("PUT")
	cluster.HandleFunc("/image", s.handleGetImages).Methods("GET")
	cluster.HandleFunc("/image/{image_id}", s.handleGetImageByID).Methods("GET")
	cluster.HandleFunc("/keyPair", s.handleGetKeyPairs).Methods("GET")
	cluster.HandleFunc("/keyPair/{keyPair_id}", s.handleGetKeyPairByID).Methods("GET")
	cluster.HandleFunc("/keyPair/{keyPair_id}", s.handleDeleteKeyPair).Methods("DELETE")
	cluster.HandleFunc("/keyPair}", s.handleCreateKeyPair).Methods("POST")
	cluster.HandleFunc("/network", s.handleGetNetworks).Methods("GET")
	cluster.HandleFunc("/network/{network_id}", s.handleGetNetworkByID).Methods("GET")
	cluster.HandleFunc("/network/private", s.handleCreatePrivateNetwork).Methods("POST")
	cluster.HandleFunc("/network/private/{network_id}", s.handleDeletePrivateNetwork).Methods("DELETE")
	cluster.HandleFunc("/project/{project_id}", s.handleGetProjectByID).Methods("GET")
	cluster.HandleFunc("/role", s.handleGetRoles).Methods("GET")
	cluster.HandleFunc("/role/{role_id}", s.handleGetRoleByID).Methods("GET")
	cluster.HandleFunc("/router", s.handleGetRouters).Methods("GET")
	cluster.HandleFunc("/router/{router_id}", s.handleGetRouterByID).Methods("GET")
	cluster.HandleFunc("/router", s.handleCreateRouter).Methods("POST")
	cluster.HandleFunc("/router/add_subnet", s.handleAddSubnetToRouter).Methods("PUT")
	cluster.HandleFunc("/router/remove_subnet", s.handleRemoveRouterSubnet).Methods("PUT")
	cluster.HandleFunc("/router/remove_external_gateway", s.handleRemoveRouterGateway).Methods("HEAD")
	cluster.HandleFunc("/router/{router_id}", s.handleDeleteRouter).Methods("DELETE")
	cluster.HandleFunc("/securityGroup", s.handleCreateSecurityGroup).Methods("POST")
	cluster.HandleFunc("/securityGroup/{securityGroup_id}", s.handleDeleteSecurityGroup).Methods("DELETE")
	cluster.HandleFunc("/securityGroup", s.handleGetSecurityGroups).Methods("GET")
	cluster.HandleFunc("/securityGroup/{securityGroup_id}", s.handleGetSecurityGroupByID).Methods("GET")
	cluster.HandleFunc("/securityRule", s.handleGetSecurityRules).Methods("GET")
	cluster.HandleFunc("/securityRule/{securityRule_id}", s.handleGetSecurityRuleByID).Methods("GET")
	cluster.HandleFunc("/securityRule", s.handleCreateSecurityRule).Methods("POST")
	cluster.HandleFunc("/securityRule/{securityRule_id}", s.handleDeleteSecurityRule).Methods("DELETE")
	cluster.HandleFunc("/server", s.handleGetServers).Methods("GET")
	cluster.HandleFunc("/server/{server_id}", s.handleGetServerByID).Methods("GET")
	cluster.HandleFunc("/server", s.handleCreateServer).Methods("POST")
	cluster.HandleFunc("/server/{server_id}", s.handleDeleteServer).Methods("DELETE")
	cluster.HandleFunc("/server/{server_id}/start", s.handleStartServer).Methods("HEAD")
	cluster.HandleFunc("/server/{server_id}/stop", s.handleStopServer).Methods("HEAD")
	cluster.HandleFunc("/server/{server_id}/attach_volume", s.handleAttachVolToServer).Methods("PUT")
	cluster.HandleFunc("/subnet", s.handleGetSubnets).Methods("GET")
	cluster.HandleFunc("/subnet/{subnet_id}", s.handleGetSubnetByID).Methods("GET")
	cluster.HandleFunc("/subnet", s.handleCreateSubnet).Methods("POST")
	cluster.HandleFunc("/subnet/{subnet_id}", s.handleDeleteSubnet).Methods("DELETE")
	cluster.HandleFunc("/user/{user_id}", s.handleGetClusterUserByID).Methods("GET")
	cluster.HandleFunc("/volume", s.handleGetVolumes).Methods("GET")
	cluster.HandleFunc("/volume/{volume_id}", s.handleGetVolumeByID).Methods("GET")
	cluster.HandleFunc("/volume/{volume_id}", s.handleDeleteVolume).Methods("DELETE")
	cluster.HandleFunc("/volume", s.handleCreateVolume).Methods("POST")
	cluster.HandleFunc("/port", s.handleGetPorts).Methods("GET")
	cluster.HandleFunc("/port/{port_id}", s.handleGetPortByID).Methods("GET")
	cluster.HandleFunc("/network/{network_id}/port", s.handleGetPortsByNetwork).Methods("GET")
	cluster.HandleFunc("/router/{device_id}/port", s.handleGetPortsByDevice).Methods("GET")
	cluster.HandleFunc("/server/{device_id}/port", s.handleGetPortsByDevice).Methods("GET")

	admin := private.PathPrefix("").Subrouter()

	admin.Use(s.isAdmin)

	admin.HandleFunc("/account", s.handleGetAllAccounts).Methods("GET")
	admin.HandleFunc("/cluster", s.handleCreateCluster).Methods("POST")
	admin.HandleFunc("/cluster/{cluster_id}", s.handleDeleteCluster).Methods("DELETE")

	clusterAdmin := admin.PathPrefix("/cluster/{cluster_id}").Subrouter()

	clusterAdmin.Use(s.authenticateClusterUser)

	clusterAdmin.HandleFunc("/flavor", s.handleCreateFlavor).Methods("POST")
	clusterAdmin.HandleFunc("/flavor/{flavor_id}", s.handleDeleteFlavor).Methods("DELETE")
	clusterAdmin.HandleFunc("/floatingIp", s.handleCreateFloatingIp).Methods("POST")
	clusterAdmin.HandleFunc("/floatingIp/{floatingIp_id}", s.handleDeleteFloatingIp).Methods("DELETE")
	clusterAdmin.HandleFunc("/image", s.handleCreateImage).Methods("POST")
	clusterAdmin.HandleFunc("/image/{image_id}/file", s.handleUploadImage).Methods("POST")
	clusterAdmin.HandleFunc("/image/{image_id}", s.handleDeleteImage).Methods("DELETE")
	clusterAdmin.HandleFunc("/network/public", s.handleCreatePublicNetwork).Methods("POST")
	clusterAdmin.HandleFunc("/network/public/{network_id}", s.handleDeletePublicNetwork).Methods("DELETE")
	clusterAdmin.HandleFunc("/project", s.handleGetProjects).Methods("GET")
	clusterAdmin.HandleFunc("/project", s.handleCreateProject).Methods("POST")
	clusterAdmin.HandleFunc("/project/{project_id}", s.handleDeleteProject).Methods("DELETE")
	clusterAdmin.HandleFunc("/role", s.handleCreateRole).Methods("POST")
	clusterAdmin.HandleFunc("/role/{role_id}", s.handleDeleteRole).Methods("DELETE")
	clusterAdmin.HandleFunc("/user", s.handleCreateClusterUser).Methods("POST")
	clusterAdmin.HandleFunc("/user", s.handleGetAllClusterUsers).Methods("GET")
	clusterAdmin.HandleFunc("/user/{user_id}", s.handleDeleteClusterUser).Methods("DELETE")
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
		return err

	case http.StatusUnprocessableEntity:
		return err

	}
	return err
}
