package models

// Book model
type Book struct {
	BookKey   string `json:"key"`
	Title        string    `json:"name"`
	EditionCount string    `json:"edition_count"`
	CoverID       float64   `json:"cover_id"`
	CoverEditionKey    string    `json:"name"`
}
