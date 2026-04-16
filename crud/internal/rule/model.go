package rule

import "gorm.io/datatypes"

type Rule struct {
	RuleID    uint           `gorm:"column:rule_id;primaryKey"`
	RuleName  string         `gorm:"column:rule_name;size:100;not null"`
	Condition datatypes.JSON `gorm:"column:condition;not null"`
	Effect    bool           `gorm:"column:effect;not null;default:false"`
}

func (Rule) TableName() string {
	return "rules"
}
