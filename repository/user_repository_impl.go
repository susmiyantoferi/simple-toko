package repository

import (
	"context"
	"errors"
	"fmt"
	"simple-toko/entity"

	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	Db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) *userRepositoryImpl {
	return &userRepositoryImpl{Db: db}
}

var (
	ErrorIdNotFound    = errors.New("id not found")
	ErrorEmailNotFound = errors.New("email not found")
	ErrorEmailExist    = errors.New("email already exist")
	ErrNotEnoughStock  = errors.New("not enough stock")
	ErrorValidation    = errors.New("validation failed")
)

func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := r.Db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrorEmailExist
		}
		return nil, fmt.Errorf("user repo: create: %w", err)
	}

	return user, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, userId uint, user *entity.User) (*entity.User, error) {
	var dataUser entity.User
	if err := r.Db.WithContext(ctx).First(&dataUser, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorIdNotFound
		}

		return nil, err
	}

	update := map[string]interface{}{}

	if user.Name != "" {
		update["name"] = user.Name
	}

	if user.Password != "" {
		update["password"] = user.Password
	}

	result := r.Db.WithContext(ctx).Model(&dataUser).Updates(update).Error
	if result != nil {
		return nil, fmt.Errorf("user repo: update: %w", result)
	}

	return &dataUser, nil

}

func (r *userRepositoryImpl) Delete(ctx context.Context, userId uint) error {
	user := entity.User{}
	result := r.Db.WithContext(ctx).Delete(&user, userId)
	if result.Error != nil {
		return fmt.Errorf("user repo: delete: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrorIdNotFound
	}

	return nil
}

func (r *userRepositoryImpl) FindById(ctx context.Context, userId uint) (*entity.User, error) {
	user := entity.User{}
	if err := r.Db.WithContext(ctx).First(&user, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorIdNotFound
		}

		return nil, fmt.Errorf("user repo: find id: %w", err)
	}

	return &user, nil
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := entity.User{}
	if err := r.Db.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorEmailNotFound
		}

		return nil, fmt.Errorf("user repo: find email: %w", err)
	}

	return &user, nil
}

func (r *userRepositoryImpl) FindAll(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error) {
	var user []*entity.User
	var totalItems int64

	if err := r.Db.WithContext(ctx).Model(&entity.User{}).Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := r.Db.WithContext(ctx).Limit(pageSize).Offset(offset).Find(&user).Error; err != nil {
		return nil, 0, err
	}

	return user, totalItems, nil
}
