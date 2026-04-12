package decision

func ToDTO(d Decision) DecisionDTO {
	return DecisionDTO{
		DecisionID: d.DecisionID,
		RequestID:  d.RequestID,
		PolicyID:   d.PolicyID,
		Result:     d.Result,
	}
}
