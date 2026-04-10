package request

func ToDTO(r Request) RequestDTO {
	return RequestDTO{
		ID:        r.RequestID,
		UserID:    r.UserID,
		HostID:    r.HostID,
		ServiceID: r.ServiceID,
		ActionID:  r.ActionID,
	}
}

func FromCreateDTO(dto CreateRequestDTO) Request {
	return Request{
		UserID:    dto.UserID,
		HostID:    dto.HostID,
		ServiceID: dto.ServiceID,
		ActionID:  dto.ActionID,
	}
}
