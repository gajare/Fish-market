package service

import (
	"context"
	"errors"
	"strings"

	"github.com/gajare/Fish-market/db"
	"github.com/gajare/Fish-market/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{db: db.DB}
}

func (s *UserService) hashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}
func (s *UserService) comparedPassword(hash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}

func (s *UserService) Create(ctx context.Context, dto models.CreateUserDTO, canSetRole bool) (models.User, error) {
	if dto.Password == "" {
		return models.User{}, errors.New("Password required")
	}
	role := models.RoleCustomer

	if canSetRole && dto.Role != nil {
		role = *dto.Role
	}

	hash, err := s.hashPassword(dto.Password)
	if err != nil {
		return models.User{}, err
	}

	u := models.User{
		FullName:     dto.FullName,
		Email:        strings.ToLower(dto.Email),
		PasswordHash: hash,
		Role:         role,
		Phone:        dto.Phone,
		Address:      dto.Address,
	}
	if err := s.db.WithContext(ctx).Create(&u).Error; err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (models.User, error) {
	var u models.User
	if err := s.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, err
		}
	}
	return u, nil
}

func (s *UserService) List(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := s.db.WithContext(ctx).Order("id desc").Find(&users).Error
	return users, err
}

func (s *UserService) Update(ctx context.Context, id uint, dto models.UpdateUserDTO, canSetRole bool) (models.User, error) {
	var u models.User
	if err := s.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	if dto.FullName != nil {
		u.FullName = *dto.FullName
	}
	if dto.Password != nil && *dto.Password != "" {
		hash, err := s.hashPassword(*dto.Password)
		if err != nil {
			return u, err
		}
		u.PasswordHash = hash
	}
	if dto.Phone != nil {
		u.Phone = dto.Phone
	}
	if dto.Address != nil {
		u.Address = dto.Address
	}
	if canSetRole && dto.Role != nil {
		u.Role = *dto.Role
	}
	if err := s.db.WithContext(ctx).Save(&u).Error; err != nil {
		return u, err
	}
	return u, nil
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (s *UserService) Login(ctx context.Context, email, password string) (models.User, error) {
	var u models.User
	if err := s.db.WithContext(ctx).Where("email=?", strings.ToLower(email)).First(&u).Error; err != nil {
		return models.User{}, errors.New("invalid credentials")
	}
	if err := s.comparedPassword(u.PasswordHash, password); err != nil {
		return models.User{}, errors.New("invalid credentils")
	}
	return u, nil
}
