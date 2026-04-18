package request

type RequestDTO struct {
	RequestID uint        `json:"request_id"`
	UserID    *uint       `json:"user_id"`
	HostID    *uint       `json:"host_id"`
	ServiceID *uint       `json:"service_id"`
	ActionID  *uint       `json:"action_id"`
	Context   interface{} `json:"context"`
}

type CreateRequestDTO struct {
	UserID    *uint       `json:"user_id"`
	HostID    *uint       `json:"host_id"`
	ServiceID *uint       `json:"service_id"`
	ActionID  *uint       `json:"action_id"`
	Context   interface{} `json:"context"`
}
