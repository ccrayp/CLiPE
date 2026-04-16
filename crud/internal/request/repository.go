package request

import "clipe/pkg/database"

type RequestRepository struct {
	db_ *database.DB
}

func NewRequestRep(db *database.DB) *RequestRepository {
	return &RequestRepository{db_: db}
}

func (r *RequestRepository) Select(filter *RequestDTO, limit int, offset int) ([]RequestDTO, error) {

	var requests []Request

	if err := r.db_.Conn().
		Limit(limit).
		Offset(offset).
		Where(filter).
		Find(&requests).Error; err != nil {
		return nil, err
	}

	var result []RequestDTO
	for _, req := range requests {
		result = append(result, ToDTO(req))
	}

	return result, nil
}

func (r *RequestRepository) Create(dto *CreateRequestDTO) (*uint, error) {

	model := FromCreateDTO(*dto)

	if err := r.db_.Conn().Create(&model).Error; err != nil {
		return nil, err
	}

	return &model.RequestID, nil
}

func (r *RequestRepository) Update(id uint, dto *CreateRequestDTO) error {

	var model Request

	if err := r.db_.Conn().First(&model, id).Error; err != nil {
		return err
	}

	model.UserID = dto.UserID
	model.HostID = dto.HostID
	model.ServiceID = dto.ServiceID
	model.ActionID = dto.ActionID

	return r.db_.Conn().Save(&model).Error
}

func (r *RequestRepository) Delete(id uint) error {
	return r.db_.Conn().Unscoped().Delete(&Request{}, id).Error
}
