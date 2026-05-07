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

	if err := r.db_.Conn().
		Limit(limit).
		Offset(offset).
		Where(&User{
			UserID:   filter.UserID,
			UserName: filter.UserName,
			UID:      filter.UID,
			GID:      filter.GID,
			HostId:   filter.HostId,
		}).
		Find(&users).Error; err != nil {
		return nil, err
	}

	var result []UserDTO
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
