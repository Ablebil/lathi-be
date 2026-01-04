package contract

import (
	"context"

	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/google/uuid"
)

type StoryUsecaseItf interface {
	GetChapterList(ctx context.Context, userID uuid.UUID) ([]dto.ChapterListReponse, *response.APIError)
	GetChapterContent(ctx context.Context, userID, chapterID uuid.UUID) (*dto.ChapterContentResponse, *response.APIError)
	GetUserSession(ctx context.Context, userID, chapterID uuid.UUID) (*dto.UserSessionResponse, *response.APIError)
	StartSession(ctx context.Context, userID, chapterID uuid.UUID) *response.APIError
	SubmitAction(ctx context.Context, userID uuid.UUID, req *dto.StoryActionRequest) (*dto.StoryActionResponse, *response.APIError)
}

type StoryRepositoryItf interface {
	GetAllChapters(ctx context.Context) ([]entity.Chapter, error)
	GetChapterByID(ctx context.Context, id uuid.UUID) (*entity.Chapter, error)
	GetSlideByID(ctx context.Context, id uuid.UUID) (*entity.Slide, error)
	FindSession(ctx context.Context, userID, chapterID uuid.UUID) (*entity.UserStorySession, error)
	CreateSession(ctx context.Context, session *entity.UserStorySession) error
	UpdateSession(ctx context.Context, session *entity.UserStorySession) error
	UnlockVocabularies(ctx context.Context, userID uuid.UUID, vocabIDs []uuid.UUID) (int64, error)
	CountChapters(ctx context.Context) (int64, error)
}
