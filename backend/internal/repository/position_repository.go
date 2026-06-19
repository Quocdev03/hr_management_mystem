package repository

import (
	"chiquoc_hocgolang/internal/model"

	"gorm.io/gorm"
)

type PositionRepository interface {
	WithTx(tx *gorm.DB) PositionRepository
	Create(pos *model.Position) error
	FindByID(id uint) (*model.Position, error)
	FindAll() ([]model.Position, error)
	UpdateFields(id uint, fields map[string]interface{}) error
	Delete(id uint) error
	ExistsByName(name string) (bool, error)
	CountEmployees(id uint) (int64, error)
}

type positionRepository struct {
	db *gorm.DB
}

func NewPositionRepository(db *gorm.DB) PositionRepository {
	return &positionRepository{db: db}
}

func (r *positionRepository) WithTx(tx *gorm.DB) PositionRepository {
	return &positionRepository{db: tx}
}

func (r *positionRepository) Create(pos *model.Position) error {
	return r.db.Create(pos).Error
}

func (r *positionRepository) FindByID(id uint) (*model.Position, error) {
	var pos model.Position
	if err := r.db.First(&pos, id).Error; err != nil {
		return nil, err
	}
	return &pos, nil
}

func (r *positionRepository) FindAll() ([]model.Position, error) {
	var positions []model.Position
	if err := r.db.Order("created_at DESC").Find(&positions).Error; err != nil {
		return nil, err
	}
	return positions, nil
}

func (r *positionRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	return r.db.Model(&model.Position{}).Where("id = ?", id).Updates(fields).Error
}

func (r *positionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Position{}, id).Error
}

func (r *positionRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Position{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func (r *positionRepository) CountEmployees(id uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Employee{}).Where("position_id = ?", id).Count(&count).Error
	return count, err
}
