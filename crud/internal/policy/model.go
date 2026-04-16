package policy

type Policy struct {
	PolicyID   uint   `gorm:"column:policy_id;primaryKey"`
	PolicyName string `gorm:"column:policy_name;size:100;not null;unique"`

	UserID    *uint `gorm:"column:user_id;not null"`
	HostID    *uint `gorm:"column:host_id;not null"`
	ServiceID *uint `gorm:"column:service_id;not null"`
	ActionID  *uint `gorm:"column:action_id;not null"`
	RuleID    *uint `gorm:"column:rule_id;not null"`

	Status bool `gorm:"column:status;not null;default:false;not null"`
}

func (Policy) TableName() string {
	return "policies"
}
