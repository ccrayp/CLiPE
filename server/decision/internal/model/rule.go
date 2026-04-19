package model

import "encoding/json"

type PolicyMatchResponse struct {
	Policy PolicyResponse `json:"policy"`
	Rule   Rule           `json:"rule"`
	Effect bool           `json:"effect"`
}

type PolicyResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	UserID uint   `json:"user_id,omitempty"`
	Status bool   `json:"status"`
}

type Condition struct {
	Type     string      `json:"type"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type Rule struct {
	Conditions []Condition `json:"conditions"`
	Effect     bool        `json:"effect"`
}

func ParseRule(raw any) (*Rule, error) {
	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	var rule Rule
	if err := json.Unmarshal(bytes, &rule); err != nil {
		return nil, err
	}

	return &rule, nil
}
