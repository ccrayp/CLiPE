package policy

import (
	"clipe/pkg/database"
)

type PolicyRepository struct {
	db_ *database.DB
}

func NewPolicyRep(db *database.DB) *PolicyRepository {
	return &PolicyRepository{
		db_: db,
	}
}

func (r *PolicyRepository) Select(policy *PolicyDTO, limit int, offset int) ([]PolicyDTO, error) {
	var policies []Policy

	if err := r.db_.Conn().Limit(limit).Offset(offset).Where(policy).Find(&policies).Error; err != nil {
		return nil, err
	}

	var result []PolicyDTO
	for _, a := range policies {
		result = append(result, ToDTO(a))
	}

	return result, nil
}

func (r *PolicyRepository) Create(policy *CreatePolicyDTO) (*uint, error) {

	model := Policy{
		PolicyName: policy.PolicyName,
		UserID:     policy.UserID,
		HostID:     policy.HostID,
		ServiceID:  policy.ServiceID,
		ActionID:   policy.ActionID,
		RuleID:     policy.RuleID,
		Status:     policy.Status,
	}

	result := r.db_.Conn().Create(&model)

	if result.Error != nil {
		return nil, result.Error
	}

	return &model.PolicyID, nil
}

func (r *PolicyRepository) Update(id uint, policy *CreatePolicyDTO) error {
	model := Policy{
		PolicyID: id,
	}

	r.db_.Conn().First(&model)

	model.PolicyName = policy.PolicyName
	model.UserID = policy.UserID
	model.HostID = policy.HostID
	model.ActionID = policy.ActionID
	model.RuleID = policy.RuleID
	model.Status = policy.Status

	result := r.db_.Conn().Save(&model)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *PolicyRepository) Delete(id uint) error {
	result := r.db_.Conn().Unscoped().Delete(&Policy{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
