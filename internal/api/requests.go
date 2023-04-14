package api

type CreateAccountRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateServerRequest struct {
	AvailabilityZone string   `json:"availability_zone"`
	PrivateIp        string   `json:"private_ip"`
	PublicIp         string   `json:"public_ip"`
	ImageID          string   `json:"image_id"`
	KeyID            string   `json:"key_id"`
	Name             string   `json:"name"`
	SecurityGroups   []string `json:"security_groups"`
	VolumesAttached  []string `json:"volumes_attached"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateClusterRequest struct {
	Location string                `json:"location"`
	Url      string                `json:"url"`
	Admin    *CreateAccountRequest `json:"admin"`
}
