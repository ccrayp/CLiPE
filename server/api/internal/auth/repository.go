package auth

import (
	"clipe/pkg/database"

	"gorm.io/gorm"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveRefreshToken(token, username string) error {
	return r.db.Conn().Create(&RefreshToken{
		Token:    token,
		Username: username,
	}).Error
}

func (r *Repository) GetRefreshToken(token string) (*RefreshToken, error) {
	var rt RefreshToken
	if err := r.db.Conn().Where("token = ?", token).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *Repository) DeleteRefreshToken(token string) error {
	return r.db.Conn().Where("token = ?", token).Delete(&RefreshToken{}).Error
}

func (r *Repository) GetUserByUsername(username string) (*UserAuthDTO, error) {
	var user UserAuthDTO

	err := r.db.Conn().
		Table("sys_users").
		Select("username, password_hash as password").
		Where("username = ?", username).
		Scan(&user).Error

	if err != nil {
		return nil, err
	}

	if user.Username == "" {
		return nil, gorm.ErrRecordNotFound
	}

	return &user, nil
}
