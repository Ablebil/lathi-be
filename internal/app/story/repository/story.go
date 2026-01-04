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

type storyRepository struct {
	db *gorm.DB
}

func NewStoryRepository(db *gorm.DB) contract.StoryRepositoryItf {
	return &storyRepository{
		db: db,
	}
}

func (r *storyRepository) GetAllChapters(ctx context.Context) ([]entity.Chapter, error) {
	var chapters []entity.Chapter
	err := r.db.WithContext(ctx).Order("order_index ASC").Find(&chapters).Error
	if err != nil {
		return nil, err
	}
	return chapters, nil
}

func (r *storyRepository) GetChapterByID(ctx context.Context, id uuid.UUID) (*entity.Chapter, error) {
	var chapter entity.Chapter
	err := r.db.WithContext(ctx).
		Preload("Slides", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Preload("Slides.Vocabularies").
		Where("id = ?", id).
		First(&chapter).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &chapter, nil
}

func (r *storyRepository) GetSlideByID(ctx context.Context, id uuid.UUID) (*entity.Slide, error) {
	var slide entity.Slide
	err := r.db.WithContext(ctx).
		Preload("Vocabularies").
		Where("id = ?", id).
		First(&slide).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &slide, nil
}

func (r *storyRepository) FindSession(ctx context.Context, userID, chapterID uuid.UUID) (*entity.UserStorySession, error) {
	var session entity.UserStorySession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND chapter_id = ?", userID, chapterID).
		First(&session).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *storyRepository) CreateSession(ctx context.Context, session *entity.UserStorySession) error {
	// upsert session
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "chapter_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"current_slide_id", "current_hearts", "is_game_over", "is_completed", "history_log", "updated_at"}),
	}).Create(session).Error
}

func (r *storyRepository) UpdateSession(ctx context.Context, session *entity.UserStorySession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *storyRepository) UnlockVocabularies(ctx context.Context, userID uuid.UUID, vocabIDs []uuid.UUID) (int64, error) {
	if len(vocabIDs) == 0 {
		return 0, nil
	}

	userVocabs := make([]entity.UserVocabulary, len(vocabIDs))
	for i, vid := range vocabIDs {
		userVocabs[i] = entity.UserVocabulary{
			UserID:       userID,
			DictionaryID: vid,
			UnlockedAt:   time.Now(),
		}
	}

	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&userVocabs)

	return result.RowsAffected, result.Error
}

func (r *storyRepository) CountChapters(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Chapter{}).Count(&count).Error
	return count, err
}
