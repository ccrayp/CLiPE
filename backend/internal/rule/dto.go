package rule

type RuleDTO struct {
	RuleID    uint        `json:"rule_id"`
	RuleName  string      `json:"rule_name"`
	Condition interface{} `json:"condition"`
	Effect    bool        `json:"effect"`
}

type CreateRuleDTO struct {
	RuleName  string      `json:"rule_name" binding:"required"`
	Condition interface{} `json:"condition" binding:"required"`
	Effect    bool        `json:"effect"`
}
