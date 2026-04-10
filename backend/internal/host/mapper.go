package host

func ToDTO(h Host) HostDTO {
	return HostDTO{
		ID: h.HostID,
		IP: h.IP,
	}
}

func FromCreateDTO(dto CreateHostDTO) Host {
	return Host{
		IP: dto.IP,
	}
}
