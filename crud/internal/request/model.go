package request

import "gorm.io/datatypes"

type Request struct {
	RequestID uint `gorm:"column:request_id;primaryKey"`

	UserID    *uint          `gorm:"column:user_id"`
	HostID    *uint          `gorm:"column:host_id"`
	ServiceID *uint          `gorm:"column:service_id"`
	ActionID  *uint          `gorm:"column:action_id"`
	Context   datatypes.JSON `gorm:"column:context;not null"`
}

func (Request) TableName() string {
	return "requests"
}
