package host

func ToDTO(h Host) HostDTO {
	return HostDTO{
		HostID: h.HostID,
		IP:     h.IP,
	}
}

func FromCreateDTO(dto CreateHostDTO) Host {
	return Host{
		IP: dto.IP,
	}
}
