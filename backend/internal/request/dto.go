package request

type RequestDTO struct {
	ID uint `json:"id"`

	UserID    uint `json:"user_id"`
	HostID    uint `json:"host_id"`
	ServiceID uint `json:"service_id"`
	ActionID  uint `json:"action_id"`
}

type CreateRequestDTO struct {
	UserID    uint `json:"user_id" binding:"required"`
	HostID    uint `json:"host_id" binding:"required"`
	ServiceID uint `json:"service_id" binding:"required"`
	ActionID  uint `json:"action_id" binding:"required"`
}
