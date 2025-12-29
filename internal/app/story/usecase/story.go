package usecase

import (
	"context"
	"encoding/json"

	"github.com/Ablebil/lathi-be/internal/config"
	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/internal/infra/minio"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/google/uuid"
)

type storyUsecase struct {
	repo    contract.StoryRepositoryItf
	storage minio.MinioItf
	env     *config.Env
}

func NewStoryUsecase(storyRepo contract.StoryRepositoryItf, storage minio.MinioItf, env *config.Env) contract.StoryUsecaseItf {
	return &storyUsecase{
		repo:    storyRepo,
		storage: storage,
		env:     env,
	}
}

func (uc *storyUsecase) GetChapterList(ctx context.Context, userID uuid.UUID) ([]dto.ChapterListReponse, *response.APIError) {
	chapters, err := uc.repo.GetAllChapters(ctx)
	if err != nil {
		return nil, response.ErrInternal("failed to get chapters")
	}

	var resp []dto.ChapterListReponse
	for _, ch := range chapters {
		isLocked := false
		if ch.OrderIndex > 1 {
			isLocked = true
		}

		resp = append(resp, dto.ChapterListReponse{
			ID:            ch.ID,
			Title:         ch.Title,
			Description:   ch.Description,
			CoverImageURL: uc.storage.GetObjectURL(ch.CoverImageURL),
			OrderIndex:    ch.OrderIndex,
			IsLocked:      isLocked,
			IsCompleted:   false,
		})
	}

	return resp, nil
}

func (uc *storyUsecase) GetChapterContent(ctx context.Context, userID uuid.UUID, chapterID uuid.UUID) (*dto.ChapterContentResponse, *response.APIError) {
	chapter, err := uc.repo.GetChapterByID(ctx, chapterID)
	if err != nil {
		return nil, response.ErrInternal("failed to get chapter")
	}
	if chapter == nil {
		return nil, response.ErrNotFound("chapter not found")
	}

	var slidesResp []dto.SlideItemResponse

	for _, slide := range chapter.Slides {
		var vocabsResp []dto.VocabItemResponse
		for _, v := range slide.Vocabularies {
			vocabsResp = append(vocabsResp, dto.VocabItemResponse{
				ID:        v.ID,
				WordKrama: v.WordKrama,
				WordNgoko: v.WordNgoko,
				WordIndo:  v.WordIndo,
			})
		}

		var choicesResp []dto.ChoiceItemResponse
		if len(slide.Choices) > 0 {
			var rawChoice []struct {
				Text string `json:"text"`
			}

			if err := json.Unmarshal(slide.Choices, &rawChoice); err != nil {
				for i, rc := range rawChoice {
					choicesResp = append(choicesResp, dto.ChoiceItemResponse{
						Index: i,
						Text:  rc.Text,
					})
				}
			}
		}

		slidesResp = append(slidesResp, dto.SlideItemResponse{
			ID:                 slide.ID,
			BackgroundImageURL: uc.storage.GetObjectURL(slide.BackgroundImageURL),
			CharacterImageURL:  uc.storage.GetObjectURL(slide.CharacterImageURL),
			AudioFileURL:       uc.storage.GetObjectURL(slide.AudioFileURL),
			SpeakerName:        slide.SpeakerName,
			Content:            slide.Content,
			NextSlideID:        slide.NextSlideID,
			Vocabularies:       vocabsResp,
			Choices:            choicesResp,
		})
	}

	return &dto.ChapterContentResponse{
		ChapterID: chapter.ID,
		Slides:    slidesResp,
	}, nil
}
