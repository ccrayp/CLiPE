package request

import (
	"time"

	"gorm.io/datatypes"
)

type Request struct {
	RequestID uint `gorm:"column:request_id;primaryKey"`

	UserID    uint           `gorm:"column:user_id"`
	Context   datatypes.JSON `gorm:"column:context;not null"`
	Timestamp time.Time      `gorm:"timestamp"`
}

func (Request) TableName() string {
	return "requests"
}
