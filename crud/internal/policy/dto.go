package policy

type PolicyDTO struct {
	PolicyID   uint   `json:"policy_id"`
	PolicyName string `json:"policy_name"`
	Status     bool   `json:"status"`

	UserID    *uint `json:"user_id,omitempty"`
	HostID    *uint `json:"host_id,omitempty"`
	ServiceID *uint `json:"service_id,omitempty"`
	ActionID  *uint `json:"action_id,omitempty"`
	RuleID    *uint `json:"rule_id,omitempty"`
}

type CreatePolicyDTO struct {
	PolicyName string `json:"policy_name"`
	Status     bool   `json:"status"`

	UserID    *uint `json:"user_id"`
	HostID    *uint `json:"host_id"`
	ServiceID *uint `json:"service_id"`
	ActionID  *uint `json:"action_id"`
	RuleID    *uint `json:"rule_id"`
}
