package request

type RequestDTO struct {
	RequestID uint        `json:"request_id"`
	UserID    *uint       `json:"user_id"`
	Context   interface{} `json:"context"`
}

type CreateRequestDTO struct {
	UserID  *uint       `json:"user_id"`
	Context interface{} `json:"context"`
}
