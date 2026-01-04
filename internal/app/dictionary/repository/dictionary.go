package repository

import (
	"context"
	"strings"

	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type dictionaryRepository struct {
	db *gorm.DB
}

func NewDictionaryRepository(db *gorm.DB) *dictionaryRepository {
	return &dictionaryRepository{
		db: db,
	}
}

func (r *dictionaryRepository) GetDictionaries(ctx context.Context, userID uuid.UUID, search string, limit, offset int) ([]dto.DictionaryResponse, int64, error) {
	var results []dto.DictionaryResponse
	var total int64
	query := r.db.Table("dictionaries AS d").
		Select("d.id, d.word_krama, d.word_ngoko, d.word_indo, CASE WHEN uv.user_id IS NULL THEN true ELSE false END as is_locked").
		Joins("LEFT JOIN user_vocabularies uv ON d.id = uv.dictionary_id AND uv.user_id = ?", userID)

	if search != "" {
		searchLower := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"(LOWER(d.word_krama) LIKE ? OR LOWER(d.word_ngoko) LIKE ? OR LOWER(d.word_indo) LIKE ?) AND uv.user_id IS NOT NULL",
			searchLower, searchLower, searchLower,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("d.word_krama ASC").
		Limit(limit).
		Offset(offset).
		Scan(&results).Error

	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
