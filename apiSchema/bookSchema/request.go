package bookSchema

type RegistrationBook struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	AuthorID        int64  `json:"author_id"`
	Translator      string `json:"translator"`
	PublisherID     int64  `json:"publisher_id"`
	PublicationYear int    `json:"publication_year"`
	Pages           int    `json:"pages"`
	Editions        int    `json:"editions"`
	SubjectID       int64  `json:"subject_id"`
}

type GetCodeBook struct {
	ID int64 `json:"id"`
}

type PaginationBook struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}
