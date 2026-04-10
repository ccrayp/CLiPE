package decision

type DecisionDTO struct {
	ID uint `json:"id"`

	RequestID uint `json:"request_id"`
	PolicyID  uint `json:"policy_id"`

	Result bool `json:"result"`
}
