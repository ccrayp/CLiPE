package user

func ToDTO(u User) UserDTO {
	return UserDTO{
		UserID:   u.UserID,
		UserName: u.UserName,
		UID:      u.UID,
		GID:      u.GID,
		HostId:   u.HostId,
	}
}

func FromCreateDTO(dto CreateUserDTO) User {
	return User{
		UserName: dto.UserName,
		UID:      dto.UID,
		GID:      dto.GID,
		HostId:   dto.HostId,
	}
}
