package policycontent

type PolicyContentDTO struct {
	PolicyID  uint `json:"policy_id"`
	ServiceID uint `json:"service_id"`
	RuleID    uint `json:"rule_id"`
}

type CreatePolicyContentDTO struct {
	PolicyID  uint `json:"policy_id"`
	ServiceID uint `json:"service_id"`
	RuleID    uint `json:"rule_id"`
}
