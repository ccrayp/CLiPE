package request

import "time"

type RequestDTO struct {
	RequestID uint        `json:"request_id"`
	UserID    *uint       `json:"user_id"`
	Context   interface{} `json:"context"`
	Timestamp time.Time   `json:"timestamp"`
}

type CreateRequestDTO struct {
	UserID    *uint       `json:"user_id"`
	Context   interface{} `json:"context"`
	Timestamp time.Time   `json:"timestamp"`
}
