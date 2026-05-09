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

func (r *PolicyRepository) Select(filter *PolicyDTO, limit int, offset int) ([]PolicyDTO, error) {
	var policies []Policy

	query := r.db_.Conn().
		Model(&Policy{}).
		Limit(limit).
		Offset(offset)

	if filter != nil {
		if filter.PolicyID != 0 {
			query = query.Where("policy_id = ?", filter.PolicyID)
		}

		if filter.PolicyName != "" {
			query = query.Where("policy_name ILIKE ?", "%"+filter.PolicyName+"%")
		}

		if filter.UserID != 0 {
			query = query.Where("user_id = ?", filter.UserID)
		}

		if filter.Status != nil {
			query = query.Where("status = ?", *filter.Status)
		}
	}

	if err := query.Find(&policies).Error; err != nil {
		return nil, err
	}

	result := make([]PolicyDTO, 0, len(policies))
	for _, p := range policies {
		result = append(result, ToDTO(p))
	}

	return result, nil
}

func (r *PolicyRepository) Count() (*int64, error) {
	var count int64
	if err := r.db_.Conn().Model(&Policy{}).Count(&count).Error; err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *PolicyRepository) Create(policy *CreatePolicyDTO) (*uint, error) {

	model := Policy{
		PolicyName: policy.PolicyName,
		UserID:     policy.UserID,
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
