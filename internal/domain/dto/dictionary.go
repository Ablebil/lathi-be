package dto

import "github.com/google/uuid"

type DictionaryListRequest struct {
	Search string `query:"search"`
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
}

type PaginationMeta struct {
	CurrentPage  int   `json:"current_page"`
	TotalPage    int   `json:"total_page"`
	TotalItems   int64 `json:"total_items"`
	ItemsPerPage int   `json:"items_per_page"`
}

type DictionaryResponse struct {
	ID        uuid.UUID `json:"id"`
	WordKrama string    `json:"word_krama"`
	WordNgoko string    `json:"word_ngoko"`
	WordIndo  string    `json:"word_indo"`
	IsLocked  bool      `json:"is_locked"`
}

type DictionaryListResponse struct {
	Items      []DictionaryResponse `json:"items"`
	Pagination PaginationMeta       `json:"pagination"`
}
