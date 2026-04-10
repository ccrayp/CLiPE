package service

type ServiceDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CreateServiceDTO struct {
	Name string `json:"name" binding:"required"`
}
