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
	GetChapterContent(ctx context.Context, userID uuid.UUID, chapterID uuid.UUID) (*dto.ChapterContentResponse, *response.APIError)
	StartSession(ctx context.Context, userID uuid.UUID, chapterID uuid.UUID) *response.APIError
	SubmitAction(ctx context.Context, userID uuid.UUID, req *dto.StoryActionRequest) (*dto.StoryActionResponse, *response.APIError)
}

type StoryRepositoryItf interface {
	GetAllChapters(ctx context.Context) ([]entity.Chapter, error)
	GetChapterByID(ctx context.Context, id uuid.UUID) (*entity.Chapter, error)
	GetSlideByID(ctx context.Context, id uuid.UUID) (*entity.Slide, error)
	FindSession(ctx context.Context, userID, chapterID uuid.UUID) (*entity.UserStorySession, error)
	CreateSession(ctx context.Context, session *entity.UserStorySession) error
	UpdateSession(ctx context.Context, session *entity.UserStorySession) error
	GetUserLastCompletedChapter(ctx context.Context, userID uuid.UUID) (int, error)
	UpdateUserLastCompletedChapter(ctx context.Context, userID uuid.UUID, orderIndex int) error
	UnlockVocabularies(ctx context.Context, userID uuid.UUID, vocabIDs []uuid.UUID) error
}
