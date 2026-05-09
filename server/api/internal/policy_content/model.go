package policycontent

type PolicyContent struct {
	PolicyID  uint `gorm:"column:policy_id;primaryKey"`
	ServiceID uint `gorm:"column:service_id;primaryKey"`
	RuleID    uint `gorm:"column:rule_id;not null"`
}

func (PolicyContent) TableName() string {
	return "policy_contents"
}
