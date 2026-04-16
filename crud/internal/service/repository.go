package service

import "clipe/pkg/database"

type ServiceRepository struct {
	db_ *database.DB
}

func NewServiceRep(db *database.DB) *ServiceRepository {
	return &ServiceRepository{db_: db}
}

func (r *ServiceRepository) Select(filter *ServiceDTO, limit int, offset int) ([]ServiceDTO, error) {

	var services []Service

	if err := r.db_.Conn().
		Limit(limit).
		Offset(offset).
		Where(filter).
		Find(&services).Error; err != nil {
		return nil, err
	}

	var result []ServiceDTO
	for _, s := range services {
		result = append(result, ToDTO(s))
	}

	return result, nil
}

func (r *ServiceRepository) Create(dto *CreateServiceDTO) (*uint, error) {

	model := FromCreateDTO(*dto)

	if err := r.db_.Conn().Create(&model).Error; err != nil {
		return nil, err
	}

	return &model.ServiceID, nil
}

func (r *ServiceRepository) Update(id uint, dto *CreateServiceDTO) error {

	var model Service

	if err := r.db_.Conn().First(&model, id).Error; err != nil {
		return err
	}

	model.ServiceName = dto.ServiceName

	return r.db_.Conn().Save(&model).Error
}

func (r *ServiceRepository) Delete(id uint) error {
	return r.db_.Conn().Unscoped().Delete(&Service{}, id).Error
}
