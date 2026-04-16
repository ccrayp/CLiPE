package action

type Action struct {
	ActionID   uint   `gorm:"column:action_id;primaryKey"`
	ActionName string `gorm:"column:action_name;size:50;unique;not null"`
}

func (Action) TableName() string {
	return "actions"
}
