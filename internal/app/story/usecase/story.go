package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Ablebil/lathi-be/internal/config"
	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/internal/domain/entity"
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

	lastCompletedOrder, err := uc.repo.GetUserLastCompletedChapter(ctx, userID)
	if err != nil {
		return nil, response.ErrInternal("failed to get user progress")
	}

	var resp []dto.ChapterListReponse
	for _, ch := range chapters {
		isLocked := ch.OrderIndex > (lastCompletedOrder + 1)
		isCompleted := ch.OrderIndex <= lastCompletedOrder

		resp = append(resp, dto.ChapterListReponse{
			ID:            ch.ID,
			Title:         ch.Title,
			Description:   ch.Description,
			CoverImageURL: uc.storage.GetObjectURL(ch.CoverImageURL),
			OrderIndex:    ch.OrderIndex,
			IsLocked:      isLocked,
			IsCompleted:   isCompleted,
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

func (uc *storyUsecase) GetUserSession(ctx context.Context, userID, chapterID uuid.UUID) (*dto.UserSessionResponse, *response.APIError) {
	session, err := uc.repo.FindSession(ctx, userID, chapterID)
	if err != nil {
		return nil, response.ErrInternal("failed to fetch session")
	}
	if session == nil {
		return nil, nil // user hasn't played this chapter yet
	}

	return &dto.UserSessionResponse{
		SessionID:      session.ID,
		CurrentSlideID: session.CurrentSlideID,
		CurrentHearts:  session.CurrentHearts,
		IsGameOver:     session.IsGameOver,
		IsCompleted:    session.IsCompleted,
	}, nil
}

func (uc *storyUsecase) StartSession(ctx context.Context, userID uuid.UUID, chapterID uuid.UUID) *response.APIError {
	chapter, err := uc.repo.GetChapterByID(ctx, chapterID)
	if err != nil {
		return response.ErrInternal("failed to fetch chapter info")
	}
	if chapter == nil {
		return response.ErrNotFound("chapter not found")
	}
	if len(chapter.Slides) == 0 {
		return response.ErrInternal("chapter has no slides")
	}

	session := &entity.UserStorySession{
		UserID:         userID,
		ChapterID:      chapterID,
		CurrentSlideID: chapter.Slides[0].ID,
		CurrentHearts:  3,
		IsGameOver:     false,
		IsCompleted:    false,
		HistoryLog:     []byte("[]"),
	}

	if err := uc.repo.CreateSession(ctx, session); err != nil {
		return response.ErrInternal("failed to starts")
	}

	return nil
}

func (uc *storyUsecase) SubmitAction(ctx context.Context, userID uuid.UUID, req *dto.StoryActionRequest) (*dto.StoryActionResponse, *response.APIError) {
	session, err := uc.repo.FindSession(ctx, userID, req.ChapterID)
	if err != nil {
		return nil, response.ErrInternal("failed to find session")
	}
	if session == nil {
		return nil, response.ErrBadRequest("session not found, please start session first")
	}
	if session.IsGameOver || session.IsCompleted {
		return nil, response.ErrBadRequest("session already ended")
	}

	currentSlide, err := uc.repo.GetSlideByID(ctx, req.SlideID)
	if err != nil || currentSlide == nil {
		return nil, response.ErrNotFound("slide not found")
	}

	var nextSlideID *uuid.UUID = currentSlide.NextSlideID
	charExpression := "neutral"
	moodImpact := 0

	// process choice if any
	if req.ChoiceIndex != nil {
		var choices []struct {
			NextSlideID       uuid.UUID `json:"next_slide_id"`
			MoodImpact        int       `json:"mood_impact"`
			CharacterReaction string    `json:"character_reaction"`
		}

		if err := json.Unmarshal(currentSlide.Choices, &choices); err != nil {
			return nil, response.ErrInternal("failed to parse slide choices")
		}

		idx := *req.ChoiceIndex
		if idx < 0 || idx >= len(choices) {
			return nil, response.ErrBadRequest("invalid choice index")
		}

		selected := choices[idx]
		nextSlideID = &selected.NextSlideID
		moodImpact = selected.MoodImpact
		charExpression = selected.CharacterReaction
	}

	// update game state
	session.CurrentHearts += moodImpact
	isGameOver := false
	message := ""

	speaker := currentSlide.SpeakerName
	if speaker == "" {
		speaker = "Panjenenganipun"
	}

	if session.CurrentHearts <= 0 {
		session.CurrentHearts = 0
		isGameOver = true
		message = fmt.Sprintf("%s kuciwo karo omonganmu. Coba maneh ya!", speaker)
	}

	session.IsGameOver = isGameOver
	if !isGameOver && nextSlideID != nil {
		session.CurrentSlideID = *nextSlideID
	}

	isCompleted := false
	if !isGameOver && nextSlideID == nil {
		isCompleted = true
		session.IsCompleted = true
		message = "Sugeng! Sampeyan wis rampung crita iki."

		chapter, _ := uc.repo.GetChapterByID(ctx, req.ChapterID)
		if chapter != nil {
			_ = uc.repo.UpdateUserLastCompletedChapter(ctx, userID, chapter.OrderIndex)
		}
	}

	session.UpdatedAt = time.Now()
	if err := uc.repo.UpdateSession(ctx, session); err != nil {
		return nil, response.ErrInternal("failed to save progress")
	}

	if len(currentSlide.Vocabularies) > 0 {
		var vocabIDs []uuid.UUID
		for _, v := range currentSlide.Vocabularies {
			vocabIDs = append(vocabIDs, v.ID)
		}
		_ = uc.repo.UnlockVocabularies(ctx, userID, vocabIDs)
	}

	finalImageURL := ""
	if currentSlide.CharacterImageURL != "" {
		finalImageURL = uc.storage.GetObjectURL(strings.TrimSuffix(currentSlide.CharacterImageURL, ".webp") + "_" + charExpression + ".webp")
	}

	return &dto.StoryActionResponse{
		IsGameOver:        isGameOver,
		IsCompleted:       isCompleted,
		Message:           message,
		RemainingHearts:   session.CurrentHearts,
		CharacterReaction: charExpression,
		CharacterImageURL: finalImageURL,
		NextSlideID:       nextSlideID,
	}, nil
}
