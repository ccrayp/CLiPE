package aggregator

type AggregatorDTO struct {
	UserName    string `json:"user_name"`
	HostIp      string `json:"host_ip"`
	ServiceName string `json:"service_name"`
	ActionName  string `json:"action_name"`
}

type PolicyMatchResponse struct {
	Policy PolicyResponse `json:"policy"`
	Rule   Rule           `json:"rule"`
	Effect bool           `json:"effect"`
}

type PolicyResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	UserID    uint   `json:"user_id,omitempty"`
	HostID    uint   `json:"host_id,omitempty"`
	ServiceID uint   `json:"service_id"`
	ActionID  uint   `json:"action_id"`
	Status    bool   `json:"status"`
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
