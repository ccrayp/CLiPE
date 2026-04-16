package decision

import "clipe/pkg/database"

type DecisionRepository struct {
	db_ *database.DB
}

func NewDecisionRep(db *database.DB) *DecisionRepository {
	return &DecisionRepository{db_: db}
}

func (r *DecisionRepository) Select(filter *DecisionDTO, limit int, offset int) ([]DecisionDTO, error) {

	var decisions []Decision

	if err := r.db_.Conn().
		Preload("Request").
		Preload("Policy").
		Limit(limit).
		Offset(offset).
		Where(filter).
		Find(&decisions).Error; err != nil {
		return nil, err
	}

	var result []DecisionDTO
	for _, d := range decisions {
		result = append(result, ToDTO(d))
	}

	return result, nil
}

func (r *DecisionRepository) Create(dto *CreateDecisionDTO) (*uint, error) {

	model := FromCreateDTO(*dto)

	if err := r.db_.Conn().Create(&model).Error; err != nil {
		return nil, err
	}

	return &model.DecisionID, nil
}

func (r *DecisionRepository) Update(id uint, dto *CreateDecisionDTO) error {

	var model Decision

	if err := r.db_.Conn().First(&model, id).Error; err != nil {
		return err
	}

	model.RequestID = dto.RequestID
	model.PolicyID = dto.PolicyID
	model.Result = dto.Result

	return r.db_.Conn().Save(&model).Error
}

func (r *DecisionRepository) Delete(id uint) error {
	return r.db_.Conn().Unscoped().Delete(&Decision{}, id).Error
}
