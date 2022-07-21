package dto

import "github.com/dinorain/pinjembuku/pkg/utils"

type BookFindResponseDto struct {
	Meta utils.PaginationMetaDto `json:"meta"`
	Data []BookReponseDto        `json:"data"`
}

type BookReponseDto struct {
	BookKey         string       `json:"key"`
	Title           string       `json:"title"`
	EditionCount    int          `json:"edition_count"`
	CoverID         float64      `json:"cover_id"`
	CoverEditionKey string       `json:"cover_edition_key"`
	Authors         []interface{} `json:"authors,omitempty"`
}
