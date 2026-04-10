package host

type HostDTO struct {
	ID uint   `json:"id"`
	IP string `json:"ip"`
}

type CreateHostDTO struct {
	IP string `json:"ip" binding:"required"`
}
