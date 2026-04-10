package service

func ToDTO(s Service) ServiceDTO {
	return ServiceDTO{
		ID:   s.ServiceID,
		Name: s.ServiceName,
	}
}

func FromCreateDTO(dto CreateServiceDTO) Service {
	return Service{
		ServiceName: dto.Name,
	}
}
