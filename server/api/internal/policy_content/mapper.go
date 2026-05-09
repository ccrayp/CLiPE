package policycontent

func ToDTO(p PolicyContent) PolicyContentDTO {
	return PolicyContentDTO{
		PolicyID:  p.PolicyID,
		ServiceID: p.ServiceID,
		RuleID:    p.RuleID,
	}
}

func FromCreateDTO(dto CreatePolicyContentDTO) PolicyContent {
	return PolicyContent{
		PolicyID:  dto.PolicyID,
		ServiceID: dto.ServiceID,
		RuleID:    dto.RuleID,
	}
}
