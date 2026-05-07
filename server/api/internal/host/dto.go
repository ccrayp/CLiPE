package host

type HostDTO struct {
	HostID uint   `json:"host_id"`
	IP     string `json:"ip"`
}

type CreateHostDTO struct {
	IP string `json:"ip" binding:"required"`
}
