package decision

type Decision struct {
	DecisionID uint  `gorm:"column:decision_id;primaryKey"`
	RequestID  uint  `gorm:"column:request_id;not null"`
	PolicyID   *uint `gorm:"column:policy_id"`
	Result     bool  `gorm:"column:result;not null"`
}

func (Decision) TableName() string {
	return "decisions"
}
