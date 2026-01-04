package contract

import (
	"context"

	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/google/uuid"
)

type DictionaryUsecaseItf interface {
	GetDictionaryList(ctx context.Context, userID uuid.UUID, req *dto.DictionaryListRequest) (*dto.DictionaryListResponse, *response.APIError)
}

type DictionaryRepositoryItf interface {
	GetDictionaries(ctx context.Context, userID uuid.UUID, search string, limit, offset int) ([]dto.DictionaryResponse, int64, error)
}
