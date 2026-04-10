package policy

func ToDTO(p Policy) PolicyDTO {
	return PolicyDTO{
		ID:     p.PolicyID,
		Name:   p.PolicyName,
		Status: p.Status,

		UserID:    p.UserID,
		HostID:    p.HostID,
		ServiceID: p.ServiceID,
		ActionID:  p.ActionID,
		RuleID:    p.RuleID,
	}
}

func FromCreateDTO(dto CreatePolicyDTO) Policy {
	return Policy{
		PolicyName: dto.Name,
		Status:     dto.Status,

		UserID:    dto.UserID,
		HostID:    dto.HostID,
		ServiceID: dto.ServiceID,
		ActionID:  dto.ActionID,
		RuleID:    dto.RuleID,
	}
}
