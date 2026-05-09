package user

import "clipe/pkg/database"

type UserRepository struct {
	db_ *database.DB
}

func NewUserRep(db *database.DB) *UserRepository {
	return &UserRepository{db_: db}
}

func (r *UserRepository) Select(filter *UserDTO, limit int, offset int) ([]UserDTO, error) {
	var users []User

	query := r.db_.Conn().
		Model(&User{}).
		Limit(limit).
		Offset(offset)

	if filter.UserID != 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.UID != 0 {
		query = query.Where("uid = ?", filter.UID)
	}
	if filter.GID != 0 {
		query = query.Where("gid = ?", filter.GID)
	}
	if filter.HostId != 0 {
		query = query.Where("host_id = ?", filter.HostId)
	}

	if filter.UserName != "" {
		query = query.Where("user_name ILIKE ?", "%"+filter.UserName+"%")
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	result := make([]UserDTO, 0, len(users))
	for _, u := range users {
		result = append(result, ToDTO(u))
	}

	return result, nil
}

func (r *UserRepository) Count() (*int64, error) {
	var count int64
	if err := r.db_.Conn().Model(&User{}).Count(&count).Error; err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *UserRepository) Create(dto *CreateUserDTO) (*uint, error) {

	model := FromCreateDTO(*dto)

	if err := r.db_.Conn().Create(&model).Error; err != nil {
		return nil, err
	}

	return &model.UserID, nil
}

func (r *UserRepository) Update(id uint, dto *CreateUserDTO) error {

	var model User

	if err := r.db_.Conn().First(&model, id).Error; err != nil {
		return err
	}

	model.UserName = dto.UserName
	model.UID = dto.UID
	model.GID = dto.GID
	model.HostId = dto.HostId

	return r.db_.Conn().Save(&model).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db_.Conn().Unscoped().Delete(&User{}, id).Error
}
