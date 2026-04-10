package rule

type RuleDTO struct {
	ID        uint        `json:"id"`
	Name      string      `json:"name"`
	Condition interface{} `json:"condition"`
	Effect    bool        `json:"effect"`
}

type CreateRuleDTO struct {
	Name      string      `json:"name" binding:"required"`
	Condition interface{} `json:"condition" binding:"required"`
	Effect    bool        `json:"effect"`
}
