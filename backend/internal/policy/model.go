package policy

type Policy struct {
	PolicyID   uint   `gorm:"column:policy_id;primaryKey"`
	PolicyName string `gorm:"column:policy_name;size:100"`

	UserID    *uint `gorm:"column:user_id"`
	HostID    *uint `gorm:"column:host_id"`
	ServiceID *uint `gorm:"column:service_id"`
	ActionID  *uint `gorm:"column:action_id"`
	RuleID    *uint `gorm:"column:rule_id"`

	Status bool `gorm:"column:status;not null;default:false"`
}

func (Policy) TableName() string {
	return "policies"
}
