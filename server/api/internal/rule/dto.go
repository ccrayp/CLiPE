package rule

type RuleDTO struct {
	RuleID    uint        `json:"rule_id"`
	RuleName  string      `json:"rule_name"`
	Condition interface{} `json:"conditions"`
	Effect    *bool       `json:"effect,omitempty"`
}

type CreateRuleDTO struct {
	RuleName  string      `json:"rule_name" binding:"required"`
	Condition interface{} `json:"conditions" binding:"required"`
	Effect    bool        `json:"effect"`
}
