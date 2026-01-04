package usecase

import (
	"context"
	"log/slog"
	"math"

	"github.com/Ablebil/lathi-be/internal/config"
	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/google/uuid"
)

type dictionaryUsecase struct {
	repo contract.DictionaryRepositoryItf
	env  *config.Env
}

func NewDictionaryUsecase(repo contract.DictionaryRepositoryItf, env *config.Env) contract.DictionaryUsecaseItf {
	return &dictionaryUsecase{
		repo: repo,
		env:  env,
	}
}

func (uc *dictionaryUsecase) GetDictionaryList(ctx context.Context, userID uuid.UUID, req *dto.DictionaryListRequest) (*dto.DictionaryListResponse, *response.APIError) {
	page := req.Page
	if page < 0 {
		return nil, response.ErrBadRequest("Halaman ga valid")
	} else if page < 1 {
		page = 1
	}

	limit := req.Limit
	if limit < 0 {
		return nil, response.ErrBadRequest("Jumlah data per halaman ga valid")
	} else if limit < 1 {
		limit = uc.env.DefaultPageLimit
	} else if limit > uc.env.MaxPageLimit {
		limit = uc.env.MaxPageLimit
	}

	offset := (page - 1) * limit

	items, total, err := uc.repo.GetDictionaries(ctx, userID, req.Search, limit, offset)
	if err != nil {
		slog.Error("failed to get dictionaries", "error", err)
		return nil, response.ErrInternal("Coba lagi nanti ya!")
	}

	for i := range items {
		if items[i].IsLocked {
			items[i].WordKrama = "???"
			items[i].WordNgoko = "???"
			items[i].WordIndo = "???"
		}
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	pagination := dto.PaginationMeta{
		CurrentPage:  page,
		TotalPage:    totalPage,
		TotalItems:   total,
		ItemsPerPage: limit,
	}

	return &dto.DictionaryListResponse{
		Items:      items,
		Pagination: pagination,
	}, nil
}
