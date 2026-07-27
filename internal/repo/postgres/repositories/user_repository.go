package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/repo"
	pgmodels "github.com/DaniilYuz/Issue-Tracker-EPAM/internal/repo/postgres/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repo.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}

	m := pgmodels.UserModelFromDomain(user)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	var m pgmodels.UserModel

	if err := r.db.WithContext(ctx).First(&m, "id = ? ", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user %q not found", userID)
		}
		return nil, fmt.Errorf("userRepository.GetByID: %w", err)
	}
	return pgmodels.UserModelToDomain(&m), nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	m := pgmodels.UserModelFromDomain(user)

	res := r.db.WithContext(ctx).Model(m).Where("id = ?", m.ID).Updates(m)
	if res.Error != nil {
		return fmt.Errorf("userRepository: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("user %q not found", user.ID)
	}

	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&pgmodels.UserModel{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("userRepository.Delete: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("user %q not found", id)
	}
	return nil
}

func (r *userRepository) ListUsers(ctx context.Context) ([]*domain.User, error) {
	var models []pgmodels.UserModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("userRepository.List: %w", err)
	}
	users := make([]*domain.User, len(models))
	for i := range models {
		users[i] = pgmodels.UserModelToDomain(&models[i])
	}
	return users, nil
}
