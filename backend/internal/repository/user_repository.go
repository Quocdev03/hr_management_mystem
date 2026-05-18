package repository

import (
	"chiquoc_hocgolang/internal/model"

	"gorm.io/gorm"
)

// UserRepository interface cho các thao tác với bảng users
type UserRepository interface {
	Create(user *model.User) error
	FindAll(query model.PaginationQuery) ([]model.User, int64, error)
	FindByID(id uint) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	UpdateFields(id uint, fields map[string]interface{}) error
	Delete(id uint) error
}

// --- UserRepository Implementation ---
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindAll(query model.PaginationQuery) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := r.db.Model(&model.User{})
	if query.Search != "" {
		db = db.Where(
			"user_name LIKE ? OR email LIKE ?",
			"%"+query.Search+"%",
			"%"+query.Search+"%",
		)
	}

	db.Count(&total)
	offset := (query.Page - 1) * query.Limit

	err := db.Preload("Role").Offset(offset).Limit(query.Limit).Order("created_at DESC").Find(&users).Error

	return users, total, err
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.
		Preload("Role").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.
		Preload("Role").
		Where("user_name = ?", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Role").Where("email = ?  ", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
