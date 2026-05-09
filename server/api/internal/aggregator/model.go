package aggregator

type AggregatorDTO struct {
	UserName    string `json:"user_name"`
	ServiceName string `json:"service_name"`
}

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
