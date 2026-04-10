package policy

type PolicyDTO struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status bool   `json:"status"`

	UserID    *uint `json:"user_id,omitempty"`
	HostID    *uint `json:"host_id,omitempty"`
	ServiceID *uint `json:"service_id,omitempty"`
	ActionID  *uint `json:"action_id,omitempty"`
	RuleID    *uint `json:"rule_id,omitempty"`
}

type CreatePolicyDTO struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`

	UserID    *uint `json:"user_id"`
	HostID    *uint `json:"host_id"`
	ServiceID *uint `json:"service_id"`
	ActionID  *uint `json:"action_id"`
	RuleID    *uint `json:"rule_id"`
}
