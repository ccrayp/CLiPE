package host

type Host struct {
	HostID uint   `gorm:"column:host_id;primaryKey"`
	IP     string `gorm:"column:ip;type:varchar(15);unique;not null;"`
}

func (Host) TableName() string {
	return "hosts"
}
