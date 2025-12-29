package repository

import (
	"context"
	"errors"

	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
