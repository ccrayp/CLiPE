package policy

func ToDTO(p Policy) PolicyDTO {
	return PolicyDTO{
		PolicyID:   p.PolicyID,
		PolicyName: p.PolicyName,
		Status:     p.Status,
		UserID:     p.UserID,
	}
}

func FromCreateDTO(dto CreatePolicyDTO) Policy {
	return Policy{
		PolicyName: dto.PolicyName,
		Status:     dto.Status,
		UserID:     dto.UserID,
	}
}
