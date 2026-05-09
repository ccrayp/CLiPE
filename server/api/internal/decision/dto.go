package decision

import "time"

type DecisionDTO struct {
	DecisionID uint      `json:"decision_id"`
	RequestID  uint      `json:"request_id"`
	PolicyID   *uint     `json:"policy_id"`
	Result     bool      `json:"result"`
	Timestamp  time.Time `json:"timestamp"`
}

type SearchDecisionDTO struct {
	DecisionID uint  `json:"decision_id"`
	RequestID  uint  `json:"request_id"`
	PolicyID   *uint `json:"policy_id"`
	Result     *bool `json:"result"`
}

type CreateDecisionDTO struct {
	RequestID uint      `json:"request_id" binding:"required"`
	PolicyID  *uint     `json:"policy_id"`
	Result    bool      `json:"result"`
	Timestamp time.Time `json:"timestamp"`
}
