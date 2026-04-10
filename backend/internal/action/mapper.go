package action

func ToDTO(a Action) ActionDTO {
	return ActionDTO{
		ID:   a.ActionID,
		Name: a.ActionName,
	}
}

func FromCreateDTO(dto CreateActionDTO) Action {
	return Action{
		ActionName: dto.Name,
	}
}
