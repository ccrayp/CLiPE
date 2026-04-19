package decision

func ToDTO(d Decision) DecisionDTO {
	return DecisionDTO{
		DecisionID: d.DecisionID,
		RequestID:  d.RequestID,
		PolicyID:   d.PolicyID,
		Result:     d.Result,
		Timestamp:  d.Timestamp,
	}
}

func FromCreateDTO(dto CreateDecisionDTO) Decision {
	clean := func(v *uint) *uint {
		if v == nil || *v == 0 {
			return nil
		}
		return v
	}

	return Decision{
		RequestID: dto.RequestID,
		PolicyID:  clean(dto.PolicyID),
		Result:    dto.Result,
		Timestamp: dto.Timestamp,
	}
}
