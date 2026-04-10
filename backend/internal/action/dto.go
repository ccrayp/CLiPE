package action

type ActionDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CreateActionDTO struct {
	Name string `json:"name" binding:"required"`
}
