package bookSchema

type RegistrationBook struct {
	Code       int64  `json:"code"`
	Name       string `json:"name"`
	Writer     string `json:"writer"`
	Translator string `json:"translator"`
	Publisher  string `json:"publisher"`
}

type GetCodeBook struct {
	Code int64 `json:"code"`
}

type PaginationBook struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}
