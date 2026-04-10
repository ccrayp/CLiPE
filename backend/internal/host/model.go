package host

type Host struct {
	HostID uint   `gorm:"column:host_id;primaryKey"`
	IP     string `gorm:"column:ip;type:char(15);not null"`
}

func (Host) TableName() string {
	return "hosts"
}
