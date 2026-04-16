package service

type ServiceDTO struct {
	ServiceID   uint   `json:"service_id"`
	ServiceName string `json:"service_name"`
}

type CreateServiceDTO struct {
	ServiceName string `json:"service_name" binding:"required"`
}
