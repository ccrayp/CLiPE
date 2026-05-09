package policy

type PolicyDTO struct {
	PolicyID   uint   `json:"policy_id"`
	PolicyName string `json:"policy_name"`
	Status     bool   `json:"status"`
	UserID     *uint  `json:"user_id,omitempty"`
}

type CreatePolicyDTO struct {
	PolicyName string `json:"policy_name"`
	Status     bool   `json:"status"`
	UserID     *uint  `json:"user_id"`
}
