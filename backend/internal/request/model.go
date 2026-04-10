package request

type Request struct {
	RequestID uint `gorm:"column:request_id;primaryKey"`

	UserID    uint `gorm:"column:user_id;not null"`
	HostID    uint `gorm:"column:host_id;not null"`
	ServiceID uint `gorm:"column:service_id;not null"`
	ActionID  uint `gorm:"column:action_id;not null"`
}

func (Request) TableName() string {
	return "requests"
}
