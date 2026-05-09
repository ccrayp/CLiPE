package policycontent

import "clipe/pkg/database"

type PolicyContentRepository struct {
	db_ *database.DB
}

func NewPolicyContentRep(db *database.DB) *PolicyContentRepository {
	return &PolicyContentRepository{
		db_: db,
	}
}

func (r *PolicyContentRepository) Select(
	filter *PolicyContentDTO,
	limit int,
	offset int,
) ([]PolicyContentDTO, error) {

	var models []PolicyContent

	if err := r.db_.Conn().
		Where(filter).
		Limit(limit).
		Offset(offset).
		Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]PolicyContentDTO, 0, len(models))
	for _, m := range models {
		result = append(result, ToDTO(m))
	}

	return result, nil
}

func (r *PolicyContentRepository) Count() (*int64, error) {
	var count int64

	if err := r.db_.Conn().
		Model(&PolicyContent{}).
		Count(&count).Error; err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *PolicyContentRepository) Create(dto *CreatePolicyContentDTO) error {
	model := FromCreateDTO(*dto)

	if err := r.db_.Conn().Create(&model).Error; err != nil {
		return err
	}

	return nil
}

func (r *PolicyContentRepository) Update(
	policyID uint,
	serviceID uint,
	dto *CreatePolicyContentDTO,
) error {

	var model PolicyContent

	if err := r.db_.Conn().
		Where("policy_id = ? AND service_id = ?", policyID, serviceID).
		First(&model).Error; err != nil {
		return err
	}

	model.RuleID = dto.RuleID

	if err := r.db_.Conn().Save(&model).Error; err != nil {
		return err
	}

	return nil
}

func (r *PolicyContentRepository) Delete(
	policyID uint,
	serviceID uint,
) error {

	if err := r.db_.Conn().
		Where("policy_id = ? AND service_id = ?", policyID, serviceID).
		Delete(&PolicyContent{}).Error; err != nil {
		return err
	}

	return nil
}
