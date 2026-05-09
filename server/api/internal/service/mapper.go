package service

func ToDTO(s Service) ServiceDTO {
	return ServiceDTO{
		ServiceID:   s.ServiceID,
		ServiceName: s.ServiceName,
	}
}

func FromCreateDTO(dto CreateServiceDTO) Service {
	return Service{
		ServiceName: dto.ServiceName,
	}
}
