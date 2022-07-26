package models

// Book model
type Book struct {
	BookKey         string       `json:"key"`
	Title           string       `json:"title"`
	EditionCount    int          `json:"edition_count"`
	CoverID         float64      `json:"cover_id"`
	CoverEditionKey string       `json:"cover_edition_key"`
	Authors         []interface{} `json:"authors,omitempty"`
}
