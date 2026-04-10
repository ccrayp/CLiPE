package user

type UserDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	UID  int    `json:"uid"`
	GID  *int   `json:"gid,omitempty"`
}

type CreateUserDTO struct {
	Name string `json:"name"`
	UID  int    `json:"uid" binding:"required"`
	GID  *int   `json:"gid"`
}
