package request

import (
	"clipe/pkg/database"
	"encoding/json"
	"time"
)

type RequestRepository struct {
	db_ *database.DB
}

func NewRequestRep(db *database.DB) *RequestRepository {
	return &RequestRepository{db_: db}
}

func (r *RequestRepository) Select(filter *RequestDTO, limit int, offset int) ([]RequestDTO, error) {
	var requests []Request

	query := r.db_.Conn().Limit(limit).Offset(offset)

	if filter != nil {
		if filter.RequestID != 0 {
			query = query.Where("request_id = ?", filter.RequestID)
		}
		if filter.UserID != nil {
			query = query.Where("user_id = ?", *filter.UserID)
		}
	}

	if err := query.Find(&requests).Error; err != nil {
		return nil, err
	}

	result := make([]RequestDTO, 0, len(requests))
	for _, req := range requests {
		result = append(result, ToDTO(req))
	}

	return result, nil
}

func (r *RequestRepository) Create(dto *CreateRequestDTO) (*uint, error) {
	dto.Timestamp = time.Now()
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

	if dto.Context != nil {
		condBytes, err := json.Marshal(dto.Context)
		if err != nil {
			return err
		}
		model.Context = condBytes
	}

	return r.db_.Conn().Save(&model).Error
}

func (r *RequestRepository) Delete(id uint) error {
	return r.db_.Conn().Unscoped().Delete(&Request{}, id).Error
}
