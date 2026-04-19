package aggregator

import (
	"clipe/internal/policy"
	"clipe/internal/rule"
	"clipe/internal/user"
	"clipe/pkg/database"
)

type AggreagtorRepository struct {
	db *database.DB
}

func NewAggregatorRepository(db *database.DB) *AggreagtorRepository {
	return &AggreagtorRepository{
		db: db,
	}
}

func (r *AggreagtorRepository) FindUserIdByName(userName string) (*uint, error) {
	filter := user.UserDTO{
		UserName: userName,
	}

	record := user.User{}

	err := r.db.Conn().Where(filter).Find(&record).Error
	if err != nil {
		return nil, err
	}

	return &record.UserID, nil
}

func (r *AggreagtorRepository) FindPolicy(userId uint) (*policy.Policy, error) {
	filter := policy.PolicyDTO{
		UserID: &userId,
	}

	record := policy.Policy{}

	err := r.db.Conn().Where(filter).Find(&record).Error
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *AggreagtorRepository) FindRuleById(ruleId uint) (*rule.Rule, error) {
	var record rule.Rule

	err := r.db.Conn().
		Where("rule_id = ?", ruleId).
		First(&record).Error

	if err != nil {
		return nil, err
	}

	return &record, nil
}
