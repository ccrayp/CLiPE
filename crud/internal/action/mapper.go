package action

func ToDTO(a Action) ActionDTO {
	return ActionDTO{
		ActionID:   a.ActionID,
		ActionName: a.ActionName,
	}
}

func FromCreateDTO(dto CreateActionDTO) Action {
	return Action{
		ActionName: dto.ActionName,
	}
}
