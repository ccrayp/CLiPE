package decision

func ToDTO(d Decision) DecisionDTO {
	return DecisionDTO{
		ID:        d.DecisionID,
		RequestID: d.RequestID,
		PolicyID:  d.PolicyID,
		Result:    d.Result,
	}
}
