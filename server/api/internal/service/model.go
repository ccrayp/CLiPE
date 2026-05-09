package service

type Service struct {
	ServiceID   uint   `gorm:"column:service_id;primaryKey"`
	ServiceName string `gorm:"column:service_name;size:50;unique;not null"`
}

func (Service) TableName() string {
	return "services"
}
