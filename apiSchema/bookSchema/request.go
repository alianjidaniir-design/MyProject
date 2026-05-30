package bookSchema

type RegistrationBook struct {
	Name            string `json:"name" validate:"required,max=63"`
	AuthorID        int64  `json:"author_id" validate:"required"`
	TranslatorID    *int64 `json:"translator_id" validate:"omitempty"`
	PublisherID     int64  `json:"publisher_id" validate:"required"`
	PublicationYear int    `json:"publication_year" validate:"required,min=1000,max=9999"`
	Pages           int    `json:"pages" validate:"required"`
	Editions        int    `json:"edition" validate:"required"`
	SubjectID       int64  `json:"subject_id" validate:"required"`
}

type GetCodeBook struct {
	ID int64 `json:"id" validate:"required"`
}

type PaginationBook struct {
	Page         int    `json:"page" validate:"required"`
	PageSize     int    `json:"page_size" validate:"required"`
	AuthorID     *int64 `json:"author_id" validate:"omitempty"`
	TranslatorID *int64 `json:"translator_id" validate:"omitempty"`
	PublisherID  *int64 `json:"publisher_id" validate:"omitempty"`
	SubjectID    *int64 `json:"subject_id" validate:"omitempty"`
}
