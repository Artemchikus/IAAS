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

	// cluster.HandleFunc("/flavor").Methods("GET")
	cluster.HandleFunc("/flavor/{flavor_id}", s.handleGetFlavorByID).Methods("GET")
	// cluster.HandleFunc("/floatingIp").Methods("GET")
	// cluster.HandleFunc("/floatingIp/{floatingIp_id}").Methods("GET")
	// cluster.HandleFunc("/image").Methods("GET")
	// cluster.HandleFunc("/image/{image_id}").Methods("GET")
	// cluster.HandleFunc("/keyPair").Methods("GET")
	// cluster.HandleFunc("/keyPair/{keyPair_id}").Methods("GET")
	// cluster.HandleFunc("/keyPair/{keyPair_id}").Methods("DELETE")
	// cluster.HandleFunc("/keyPair}").Methods("POST")
	// cluster.HandleFunc("/network").Methods("GET")
	// cluster.HandleFunc("/network/{network_id}").Methods("GET")
	// cluster.HandleFunc("/project/{project_id}").Methods("GET")
	// cluster.HandleFunc("/role/{role_id}").Methods("GET")
	// cluster.HandleFunc("/router").Methods("GET")
	// cluster.HandleFunc("/router/{router_id}").Methods("GET")
	// cluster.HandleFunc("/securityGroup").Methods("GET")
	// cluster.HandleFunc("/securityGroup/{securityGroup_id}").Methods("GET")
	// cluster.HandleFunc("/securityRule").Methods("GET")
	// cluster.HandleFunc("/securityRule/{securityRule_id}").Methods("GET")
	// cluster.HandleFunc("/server").Methods("GET")
	// cluster.HandleFunc("/server/{server_id}").Methods("GET")
	// cluster.HandleFunc("/server/{server_id}").Methods("CREATE")
	// cluster.HandleFunc("/server/{server_id}").Methods("DELETE")
	// cluster.HandleFunc("/subnet").Methods("GET")
	// cluster.HandleFunc("/subnet/{subnet_id}").Methods("GET")
	cluster.HandleFunc("/user/{user_id}", s.handleGetClusterUserByID).Methods("GET")
	// cluster.HandleFunc("/volume").Methods("GET")
	// cluster.HandleFunc("/volume/{volume_id}").Methods("GET")
	// cluster.HandleFunc("/volume/{volume_id}").Methods("DELETE")
	// cluster.HandleFunc("/volume/{volume_id}").Methods("POST")

	admin := private.PathPrefix("").Subrouter()

	admin.Use(s.isAdmin)

	admin.HandleFunc("/account", s.handleGetAllAccounts).Methods("GET")
	admin.HandleFunc("/cluster", s.handleCreateCluster).Methods("POST")
	admin.HandleFunc("/cluster/{cluster_id}", s.handleDeleteCluster).Methods("DELETE")

	clusterAdmin := admin.PathPrefix("/cluster/{cluster_id}").Subrouter()

	clusterAdmin.Use(s.authenticateClusterUser)

	clusterAdmin.HandleFunc("/flavor", s.handleCreateFlavor).Methods("POST")
	clusterAdmin.HandleFunc("/flavor/{flavor_id}", s.handleDeleteFlavor).Methods("DELETE")
	// clusterAdmin.HandleFunc("/floatingIp/{floatingIp_id}").Methods("POST")
	// clusterAdmin.HandleFunc("/floatingIp/{floatingIp_id}").Methods("DELETE")
	// clusterAdmin.HandleFunc("/image/{image_id}").Methods("POST")
	// clusterAdmin.HandleFunc("/image/{image_id}/file").Methods("POST")
	// clusterAdmin.HandleFunc("/image/{image_id}").Methods("DELETE")
	// clusterAdmin.HandleFunc("/network/{network_id}").Methods("POST")
	// clusterAdmin.HandleFunc("/network/{network_id}").Methods("DELETE")
	// clusterAdmin.HandleFunc("/project/{project_id}").Methods("POST")
	// clusterAdmin.HandleFunc("/project/{project_id}").Methods("DELETE")
	// clusterAdmin.HandleFunc("/role/{role_id}").Methods("POST")
	// clusterAdmin.HandleFunc("/role/{role_id}").Methods("DELETE")
	// clusterAdmin.HandleFunc("/router/{router_id}").Methods("POST")
	// clusterAdmin.HandleFunc("/router/{router_id}").Methods("DELETE")
	// clusterAdmin.HandleFunc("/securityGroup/{securityGroup_id}").Methods("POST")
	// clusterAdmin.HandleFunc("/securityGroup/{securityGroup_id}").Methods("DELETE")
	// clusterAdmin.HandleFunc("/securityRule/{securityRule_id}").Methods("POST")
	// clusterAdmin.HandleFunc("/securityRule/{securityRule_id}").Methods("DELETE")
	// clusterAdmin.HandleFunc("/subnet/{subnet_id}").Methods("POST")
	// clusterAdmin.HandleFunc("/subnet/{subnet_id}").Methods("DELETE")
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

// wget https://download.cirros-cloud.net/0.6.1/cirros-0.6.1-x86_64-disk.img
// openstack image create "Cirros" --file cirros-0.6.1-x86_64-disk.img --disk-format qcow2 --container-format bare --public

// openstack project create --domain default --description "Demo Project" demo

// openstack user create --domain default --project demo --password demo demo

// openstack role create CloudUser

// openstack role add --project demo --user demo CloudUser

// openstack flavor create --id 0 --vcpus 1 --ram 300 --disk 2 m1.small

// openstack router create router01

// openstack network create private --provider-network-type geneve

// openstack subnet create private-subnet --network private --subnet-range 192.168.100.0/24 --gateway 192.168.100.1

// openstack router add subnet router01 private-subnet

// openstack network create --provider-physical-network external --provider-network-type flat --external public

// openstack subnet create public-subnet --network public --subnet-range 192.168.122.0/24 --allocation-pool start=192.168.122.200,end=192.168.122.254 --gateway 192.168.122.1 --no-dhcp

// openstack router set router01 --external-gateway public

// openstack security group create secgroup01

// openstack keypair create --public-key ~/.ssh/id_rsa.pub mykey

// openstack server create --flavor m1.small --image Cirros --security-group secgroup01 --nic net-id=private --key-name mykey Cirros

// openstack floating ip create public

// openstack server add floating ip Cirros 192.168.122.204

// openstack security group rule create --protocol icmp --ingress secgroup01
// openstack security group rule create --protocol tcp --dst-port 22:22 secgroup01

// openstack console url show Cirros
// ssh cirros@192.168.122.204

// openstack volume create --size 2 disk01

// openstack server add volume Cirros disk01
