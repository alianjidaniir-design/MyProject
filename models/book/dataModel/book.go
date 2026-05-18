package dataModel

import (
	"time"
)

type Book struct {
	ID              int64      `json:"id"`
	Name            string     `json:"name"`
	AuthorId        int64      `json:"author_id"`
	TranslatorID    *int64     `json:"translator_id"`
	PublisherID     int64      `json:"publisher_id"`
	PublicationYear int        `json:"publication_year"`
	Pages           int        `json:"pages"`
	Edition         int        `json:"edition"`
	SubjectID       int64      `json:"subject_id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}
