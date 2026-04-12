package decision

type DecisionDTO struct {
	DecisionID uint `json:"decision_id"`

	RequestID uint `json:"request_id"`
	PolicyID  uint `json:"policy_id"`

	Result bool `json:"result"`
}
