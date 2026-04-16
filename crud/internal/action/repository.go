package action

import "clipe/pkg/database"

type ActionRepository struct {
	db_ *database.DB
}

func NewActionRep(db *database.DB) *ActionRepository {
	return &ActionRepository{db_: db}
}

func (r *ActionRepository) Select(filter *ActionDTO, limit int, offset int) ([]ActionDTO, error) {
	var actions []Action

	if err := r.db_.Conn().
		Limit(limit).
		Offset(offset).
		Where(filter).
		Find(&actions).Error; err != nil {
		return nil, err
	}

	var result []ActionDTO
	for _, a := range actions {
		result = append(result, ToDTO(a))
	}

	return result, nil
}

func (r *ActionRepository) Create(dto *CreateActionDTO) (*uint, error) {

	model := FromCreateDTO(*dto)

	if err := r.db_.Conn().Create(&model).Error; err != nil {
		return nil, err
	}

	return &model.ActionID, nil
}

func (r *ActionRepository) Update(id uint, dto *CreateActionDTO) error {

	var model Action

	if err := r.db_.Conn().First(&model, id).Error; err != nil {
		return err
	}

	model.ActionName = dto.ActionName

	return r.db_.Conn().Save(&model).Error
}

func (r *ActionRepository) Delete(id uint) error {
	return r.db_.Conn().Unscoped().Delete(&Action{}, id).Error
}
