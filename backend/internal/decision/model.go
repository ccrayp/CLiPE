package decision

import (
	"clipe/internal/policy"
	"clipe/internal/request"
)

type Decision struct {
	DecisionID uint `gorm:"column:decision_id;primaryKey"`

	RequestID uint `gorm:"column:request_id;not null"`
	PolicyID  uint `gorm:"column:policy_id;not null"`

	Result bool `gorm:"column:result;not null"`

	Request request.Request `gorm:"foreignKey:RequestID"`
	Policy  policy.Policy   `gorm:"foreignKey:PolicyID"`
}

func (Decision) TableName() string {
	return "decisions"
}
