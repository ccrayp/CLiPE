package user

func ToDTO(u User) UserDTO {
	return UserDTO{
		ID:   u.UserID,
		Name: u.UserName,
		UID:  u.UID,
		GID:  u.GID,
	}
}

func FromCreateDTO(dto CreateUserDTO) User {
	return User{
		UserName: dto.Name,
		UID:      dto.UID,
		GID:      dto.GID,
	}
}
