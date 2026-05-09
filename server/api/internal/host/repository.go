package host

import "clipe/pkg/database"

type HostRepository struct {
	db_ *database.DB
}

func NewHostRep(db *database.DB) *HostRepository {
	return &HostRepository{db_: db}
}

func (r *HostRepository) Select(filter *HostDTO, limit int, offset int) ([]HostDTO, error) {
	var hosts []Host

	query := r.db_.Conn().
		Model(&Host{}).
		Limit(limit).
		Offset(offset)

	if filter != nil {
		if filter.HostID != 0 {
			query = query.Where("host_id = ?", filter.HostID)
		}

		if filter.IP != "" {
			query = query.Where("ip ILIKE ?", "%"+filter.IP+"%")
		}
	}

	if err := query.Find(&hosts).Error; err != nil {
		return nil, err
	}

	result := make([]HostDTO, 0, len(hosts))
	for _, h := range hosts {
		result = append(result, ToDTO(h))
	}

	return result, nil
}

func (r *HostRepository) Count() (*int64, error) {
	var count int64
	if err := r.db_.Conn().Model(&Host{}).Count(&count).Error; err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *HostRepository) Create(dto *CreateHostDTO) (*uint, error) {

	model := FromCreateDTO(*dto)

	if err := r.db_.Conn().Create(&model).Error; err != nil {
		return nil, err
	}

	return &model.HostID, nil
}

func (r *HostRepository) Update(id uint, dto *CreateHostDTO) error {

	var model Host

	if err := r.db_.Conn().First(&model, id).Error; err != nil {
		return err
	}

	model.IP = dto.IP

	return r.db_.Conn().Save(&model).Error
}

func (r *HostRepository) Delete(id uint) error {
	return r.db_.Conn().Unscoped().Delete(&Host{}, id).Error
}
