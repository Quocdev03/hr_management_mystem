package service

import (
	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/repository"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type PositionService interface {
	Create(req model.CreatePositionRequest) (*model.Position, error)
	GetAll() ([]model.Position, error)
	GetByID(id uint) (*model.Position, error)
	Update(id uint, req model.UpdatePositionRequest) (*model.Position, error)
	Delete(id uint) error
}

type positionService struct {
	db       *gorm.DB
	posRepo  repository.PositionRepository
}

func NewPositionService(db *gorm.DB, posRepo repository.PositionRepository) PositionService {
	return &positionService{
		db:       db,
		posRepo:  posRepo,
	}
}

func (s *positionService) Create(req model.CreatePositionRequest) (*model.Position, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, errors.New("tên chức vụ không được để trống")
	}

	// Kiểm tra trùng tên
	exists, err := s.posRepo.ExistsByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("lỗi kiểm tra trùng tên: %w", err)
	}
	if exists {
		return nil, errors.New("chức vụ đã tồn tại")
	}

	pos := &model.Position{
		Name:        req.Name,
		Description: strings.TrimSpace(req.Description),
	}

	if err := s.posRepo.Create(pos); err != nil {
		return nil, fmt.Errorf("lỗi tạo chức vụ: %w", err)
	}

	return s.posRepo.FindByID(pos.ID)
}

func (s *positionService) GetAll() ([]model.Position, error) {
	return s.posRepo.FindAll()
}

func (s *positionService) GetByID(id uint) (*model.Position, error) {
	if id == 0 {
		return nil, errors.New("id chức vụ phải lớn hơn 0")
	}
	pos, err := s.posRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("không tìm thấy chức vụ này")
		}
		return nil, err
	}
	return pos, nil
}

func (s *positionService) Update(id uint, req model.UpdatePositionRequest) (*model.Position, error) {
	if id == 0 {
		return nil, errors.New("id chức vụ phải lớn hơn 0")
	}

	pos, err := s.posRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("không tìm thấy chức vụ này")
		}
		return nil, err
	}

	updateData := map[string]interface{}{}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("tên chức vụ không được để trống")
		}
		if name != pos.Name {
			// Kiểm tra trùng tên
			exists, err := s.posRepo.ExistsByName(name)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, errors.New("chức vụ đã tồn tại")
			}
			updateData["name"] = name
		}
	}

	if req.Description != nil {
		updateData["description"] = strings.TrimSpace(*req.Description)
	}

	if len(updateData) > 0 {
		if err := s.posRepo.UpdateFields(id, updateData); err != nil {
			return nil, fmt.Errorf("lỗi cập nhật chức vụ: %w", err)
		}
	}

	return s.posRepo.FindByID(id)
}

func (s *positionService) Delete(id uint) error {
	if id == 0 {
		return errors.New("id chức vụ phải lớn hơn 0")
	}

	if _, err := s.posRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("không tìm thấy chức vụ này")
		}
		return err
	}

	// Kiểm tra xem có nhân viên nào đang giữ chức vụ này không
	count, err := s.posRepo.CountEmployees(id)
	if err != nil {
		return fmt.Errorf("lỗi kiểm tra nhân viên giữ chức vụ: %w", err)
	}
	if count > 0 {
		return errors.New("không thể xóa chức vụ đang có nhân viên đảm nhiệm")
	}

	return s.posRepo.Delete(id)
}
