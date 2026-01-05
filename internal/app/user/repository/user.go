package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) contract.UserRepositoryItf {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserWithBadges(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Preload("UserBadges.Badge").Where("id = ?", id).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) GetUserLastCompletedChapter(ctx context.Context, userID uuid.UUID) (int, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Select("last_chapter_completed").First(&user, userID).Error; err != nil {
		return 0, err
	}
	return user.LastChapterCompleted, nil
}

func (r *userRepository) UpdateUserLastCompletedChapter(ctx context.Context, userID uuid.UUID, orderIndex int) error {
	// udpate only if new orderIndex > current value
	return r.db.WithContext(ctx).Model(&entity.User{}).
		Where("id = ? AND last_chapter_completed < ?", userID, orderIndex).
		Update("last_chapter_completed", orderIndex).Error
}

func (r *userRepository) IncrementUserWordCount(ctx context.Context, userID uuid.UUID, amount int) error {
	return r.db.WithContext(ctx).Model(&entity.User{}).
		Where("id = ?", userID).
		UpdateColumn("total_words_collected", gorm.Expr("total_words_collected + ?", amount)).Error
}

func (r *userRepository) UpdateUserTitle(ctx context.Context, userID uuid.UUID, title entity.Title) error {
	return r.db.WithContext(ctx).Model(&entity.User{}).
		Where("id = ?", userID).
		Update("current_title", title).Error
}

func (r *userRepository) AssignBadge(ctx context.Context, userID uuid.UUID, badgeCode string) error {
	var badge entity.Badge
	if err := r.db.WithContext(ctx).Where("code = ?", badgeCode).First(&badge).Error; err != nil {
		return err
	}

	userBadge := entity.UserBadge{
		UserID:   userID,
		BadgeID:  badge.ID,
		EarnedAt: time.Now(),
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&userBadge).Error
}
