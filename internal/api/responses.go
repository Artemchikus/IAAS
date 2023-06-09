package api

type CreateAccountResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginResponse struct {
	JWT     string `json:"jwt-token"`
	Refresh string `json:"refresh-token"`
}

type RefreshTokenResponse struct {
	JWT     string `json:"jwt-token"`
	Refresh string `json:"refresh-token"`
}

type GetAccountResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type GetAllAccountsResponse []*GetAccountResponse

type DeleteAccountResponse struct {
	DeletedID int `json:"deleted_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateClusterResponse struct {
	ID        int    `json:"id"`
	Location  string `json:"location"`
	AdminName string `json:"admin_name"`
	URL       string `json:"url"`
}

type DeleteClusterResponse struct {
	DeletedID int `json:"deleted_id"`
}

type GetClusterResponse struct {
	ID       int    `json:"id"`
	Location string `json:"location"`
	URL      string `json:"url"`
}

type GetAllClustersResponse []*GetClusterResponse

type DeleteOpenstackResurceResponse struct {
	DeletedID string `json:"deleted_id"`
}

type StartServerResponse struct {
	StartedID string `json:"started_id"`
}

type StopServerResponse struct {
	StopedID string `json:"stoped_id"`
}

type AttachVolumeResponse struct {
	VolumeID string `json:"volume_id"`
	ServerID string `json:"server_id"`
}

type RemoveOrAddSubnetResponse struct {
	SubnetID string `json:"subnet_id"`
	RouterID string `json:"router_id"`
}

type RemoveRouterExternalGatewayResponse struct {
	RouterID string `json:"router_id"`
}

type AddFloatingIPResponse struct {
	FloatingIpID string `json:"floatingIp_id"`
	PortID       string `json:"port_id"`
}
type UploadImageResponse struct {
	ImageID string `json:"image_id"`
}