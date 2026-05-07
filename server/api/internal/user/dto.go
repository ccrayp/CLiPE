package user

type UserDTO struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	UID      int    `json:"uid"`
	GID      int    `json:"gid"`
	HostId   *int   `json:"host_id"`
}

type CreateUserDTO struct {
	UserName string `json:"user_name"`
	UID      int    `json:"uid" binding:"required"`
	GID      int    `json:"gid"`
	HostId   *int   `json:"host_id"`
}
