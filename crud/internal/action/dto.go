package action

type ActionDTO struct {
	ActionID   uint   `json:"action_id"`
	ActionName string `json:"action_name"`
}

type CreateActionDTO struct {
	ActionName string `json:"action_name" binding:"required"`
}
