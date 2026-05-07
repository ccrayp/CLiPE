package user

type User struct {
	UserID   uint   `gorm:"column:user_id;primaryKey"`
	UserName string `gorm:"column:user_name;size:100;unique;not null"`
	UID      int    `gorm:"column:uid;not null"`
	GID      int    `gorm:"column:gid"`
	HostId   *int   `gorm:"column:host_id"`
}

func (User) TableName() string {
	return "users"
}
