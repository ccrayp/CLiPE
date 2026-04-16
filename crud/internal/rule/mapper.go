package rule

import "encoding/json"

func ToDTO(r Rule) RuleDTO {
	var condition interface{}
	_ = json.Unmarshal(r.Condition, &condition)

	return RuleDTO{
		RuleID:    r.RuleID,
		RuleName:  r.RuleName,
		Condition: condition,
		Effect:    r.Effect,
	}
}

func FromCreateDTO(dto CreateRuleDTO) Rule {
	condBytes, _ := json.Marshal(dto.Condition)

	return Rule{
		RuleName:  dto.RuleName,
		Condition: condBytes,
		Effect:    dto.Effect,
	}
}
